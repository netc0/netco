package def

const (
	Mail_Heartbeat = iota
	Mail_Reg
	Mail_AddRoute
	Mail_RequestData
	Mail_ResponseData
	Mail_PushData
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
	IsRequest bool
	RequestId uint32
	Route uint32
	Data      []byte
}