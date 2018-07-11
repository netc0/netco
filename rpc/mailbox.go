package rpc

import (
	"net"
	"log"
	"io"
	"errors"
	"time"
)

type IMailBox interface {
	Start()
	Stop()
	SetHandler(handler MailHandler)
	SendTo(string, *Mail) error
	Connect(string) error
}

func NewMailBox (address string) IMailBox {
	var node xMailBox
	node.bindAddress = address
	node.routines = make(map[string]* xRoutine)
	return &node
}

// 路由
type xRoutine struct {
	conn net.Conn
	remote string
	isRunning bool
	inHeatbeat bool

	parent *xMailBox
}

type xMailBox struct {
	IMailBox
	bindAddress string
	isRunning bool
	handler MailHandler
	protocol mailProtocol

	routines map[string] *xRoutine
}

// Start
func (this *xMailBox) Start () {
	var l, err = net.Listen("tcp", this.bindAddress)
	if err != nil {
		log.Println(err)
		return
	}
	this.isRunning = true

	// heart beat
	go func() {

	}()

	for this.isRunning {
		var conn, err = l.Accept()
		if err != nil {
			break
		}
		go this.handleConnection(conn)
	}
}

// Stop
func (this *xMailBox) Stop(){
	this.isRunning = false
}

// SetHandler
func (this *xMailBox) SetHandler(handler MailHandler) {
	this.handler = handler
}
// SendTo
func (this *xMailBox) SendTo(remote string, mail *Mail) error {
	var r *xRoutine
	var err error
	if r, err = this.getRoutine(remote, true); err == nil {
		return r.Send(mail)
	}
	return err
}

// 连接到远程
func (this* xMailBox) Connect(r string) error {
	_, err := this.getRoutine(r, true)
	return err
}

// 获取路径
func (this* xMailBox) getRoutine(remote string, connect bool) (*xRoutine, error){
	if this.routines[remote] == nil { // 创建mail路径
		if connect {
			var routine xRoutine
			routine.remote = remote
			routine.parent = this
			this.routines[remote] = &routine
			if err := routine.Connect(); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("routine not exist")
		}
	}
	return this.routines[remote], nil
}

func (this *xMailBox) handleConnection(conn net.Conn) {
	defer conn.Close()
	var p mailProtocol
	for this.isRunning {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {

			}
			break
		}
		data := buf[:size]
		mail := p.handleBytes(data)
		for {
			if mail == nil {
				break
			}
			if this.handler != nil {
				this.handler.OnNewMail(*mail)
			} else {
				log.Println("没有指定 handler", mail)
			}
			mail = p.handleBytes(nil)
		}
	}
}

// 发送
func (this *xRoutine) Send(mail *Mail) error {
	var protocol mailProtocol
	data := protocol.encode(mail)
	_, err := this.conn.Write(data)
	return err
}

func (this *xRoutine) Connect() error {
	c, err := net.Dial("tcp", this.remote)
	if err != nil {
		return err
	}
	this.conn = c
	this.isRunning = true
	this.parent.handler.OnRoutineConnected(this.remote)
	go this.SendHeartBeat()
	return nil
}

func (this *xRoutine) SendHeartBeat() {
	if this.inHeatbeat {
		return
	}
	this.inHeatbeat = true

	// 一直发心跳包
	ticker := time.NewTicker(time.Second)
	m := Mail{Type:0}
	for range ticker.C {
		if !this.isRunning {
			this.inHeatbeat = false
			break
		}
		if err := this.Send(&m); err != nil {
			if e:=this.Connect(); e != nil {
				this.parent.handler.OnRoutineDisconnect(this.remote, e)
			}
		}
	}

}