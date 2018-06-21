package connector

import (
	"net"
	"log"
	"time"
)

type TCPTransporter struct{
	hostPort string
	running bool
	sessions map[string]*Session
	onNewConnection func(string, *Session)
	onCloseConnection func(string)
	onConnectionData func(*Session, uint32, uint32, []byte)
}

func (this *TCPTransporter)init() {
	this.sessions = make(map[string]*Session)
}

func CreateTCPConnector(hostPort string) *TCPTransporter {
	var connector = new (TCPTransporter)
	connector.init()
	connector.hostPort = hostPort
	return connector
}

func (this *TCPTransporter) Start(
	onNewConnection func(string, *Session),
	onCloseConnection func(string),
	onConnectionData func(*Session, uint32, uint32, []byte)) {
		// 设置回调
		this.onNewConnection = onNewConnection
		this.onCloseConnection = onCloseConnection
		this.onConnectionData = onConnectionData
	go this.start()
}

func (this *TCPTransporter) start() {
	this.running = true
	var l, err = net.Listen("tcp", this.hostPort)
	if err != nil {
		log.Println(err)
	}
	log.Println("xx启动 TCP", this.hostPort)
	defer l.Close()
	defer log.Println("Close TCP Server")
	defer this.releaseSessions()

	// heart beat service
	go func() {
		var heartBeatService = time.NewTicker(time.Second)
		for range heartBeatService.C {
			go this.checkHeartBeat()
		}
	}()

	for {
		if this.running == false {
			break;
		}
		var conn, err = l.Accept()
		if err != nil {
			break
		}
		go this.handleConnection(conn)
	}
}

func (this *TCPTransporter) handleConnection(conn net.Conn) {
	var session = NewSession(conn)
	this.sessions[session.id] = session

	defer session.onClose()
	defer this.onCloseConnection(session.id)
	defer log.Println("Close session", session.id)

	session.heartBeatTime = time.Now() // 更新心跳包
	session.ok = true
	session.transporter = this
	this.onNewConnection(session.id, session)

	for {
		if !session.ok {
			break
		}
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			log.Println("读数据错误",err)
			return
		}
		data := buf[:size]
		session.handleBytes(data)
	}
}

func (this *TCPTransporter) checkHeartBeat() {
	var die []*Session
	for _ ,session := range this.sessions {
		if time.Now().Second() - session.heartBeatTime.Second() > 5 {
			die = append(die, session)
		}
	}
	for _, session := range die{
		log.Println("session:", session.id, "失去心跳")
		// 踢下线
		session.Push(PacketType_KICK, nil)
		delete(this.sessions, session.id)
		session.conn.Close()
	}
}

func (this *TCPTransporter) releaseSessions() {
	for _,session := range this.sessions {
		session.conn.Close()
	}

}

func (this * TCPTransporter) GetSession(id string) (*Session){
	return this.sessions[id]
}