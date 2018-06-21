package connector

import (
	"net"
	"log"
	"time"
)

type Session struct {
	id            string
	conn          net.Conn
	reader        PacketReader
	heartBeatTime time.Time
	transporter   *TCPTransporter
	aUDPTransporter   *UDPTransporter
	ok             bool
}

func NewSession(conn net.Conn) *Session{
	return &Session{id:conn.RemoteAddr().String(), conn:conn}
}

func (this *Session) handleBytes(data []byte) {
	var pkg = this.reader.ParsePacket(data)
	for {
		if pkg == nil {
			break
		}
		if this.procPacket(*pkg) != 0 { // pkg error, disconnect now
			this.onClose()
			break
		}
		pkg = this.reader.ParsePacket(nil)
	}
}

func (this *Session) onClose() {
	log.Println("会话关闭:", this.id)
	this.ok = false;
	this.conn.Close()
}

func (this *Session) Response(Type int, requestId uint32, data[]byte) {
	if !this.ok {return;}
	_, err := this.conn.Write(PacketResponseToBinary(Type, requestId, data))
	if err != nil {
		log.Println(err)
		this.onClose()
	}
}

func (this *Session) Push(Type int, data[]byte) {
	if !this.ok {return;}
	r, err := this.conn.Write(PacketToBinary(Type, data))
	if err != nil {
		log.Println(err)
		this.onClose()
	}
	log.Println("发送了:", r)
}

func (this *Session) GetId() string {
	return this.id
}
