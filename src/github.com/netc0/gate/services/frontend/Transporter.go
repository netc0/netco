package frontend

type ITransporter interface {
	start()
	releaseSessions()
	checkHeartBeat()
}

type Transporter struct {
	ITransporter
	running bool   // 是否在运行中
	Host    string // 绑定的Host
	OnNewConnection func(interface{})
	OnDataPacket func(interface{}, uint32, uint32, []byte) // 收到消息
}