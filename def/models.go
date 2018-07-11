package def

const (
	Mail_Heartbeat = iota
	Mail_Reg
	Mail_AddRoute
	Mail_RequestData
	Mail_ResponseData
	Mail_PushData
	Mail_ClientLeaveNotifyMe
	Mail_ClientLeaveNotification
)

type MailOffice struct {
	Address string
	Name    string
}

type MailRoutineInfo struct {
	Name string
	Routes []uint32
}

type MailClientData struct {
	ClientId  string
	Type      uint32 // 0 是request消息
	RequestId uint32
	Route uint32
	Data      []byte
	SourceAddress    string
	SourceName    string
}