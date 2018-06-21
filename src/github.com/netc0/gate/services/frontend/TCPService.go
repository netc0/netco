package frontend

import (
	"net"
	"log"
	"time"
)


func (this *TCPTransporter) start() {
	this.running = true
	var l, err = net.Listen("tcp", this.Host)
	if err != nil {
		log.Println(err)
		this.running = false
	}
	log.Println("Frontend启动 TCP", this.Host)
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

func (this *TCPTransporter) releaseSessions(){
	ClearSession()
}

func (this *TCPTransporter) checkHeartBeat() {
	var die []ISession
	ForeachSession(func(s ISession) {
		if s.IsTimeout() {
			die = append(die, s)
		}
	})

	for _, s := range die{
		log.Println("session:", s, "失去心跳")
		s.Kick()  // 踢下线
		s.Close() //关闭
	}
}

func (this *TCPTransporter) handleConnection(conn net.Conn) {
	var session TCPSession

	defer conn.Close()
	defer RemoveSession(session)

	session.OnDataPacket = this.OnDataPacket
	session.time = time.Now() // 更新心跳
	session.isOk = true
	session.conn = conn
	session.holder = session
	session.id = conn.RemoteAddr().String()
	AddSession(&session)       // 新增会话

	//var session = NewSession(new(TCPSession))
	//var session = NewSession(conn)
	//this.sessions[session.id] = session
	//
	//defer session.onClose()
	//defer this.onCloseConnection(session.id)
	//defer log.Println("Close session", session.id)
	//
	//session.heartBeatTime = time.Now() // 更新心跳包
	//session.ok = true
	//session.transporter = this
	//this.onNewConnection(session.id, session)
	//

	for {
		if !session.IsOk() {
			log.Println("break")
			break
		}
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			log.Println("读数据错误",err)
			break
		}
		data := buf[:size]

		session.HandleBytes(data)
	}
}
