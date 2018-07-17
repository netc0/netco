package def

import "errors"

const (
	Mail_Heartbeat = iota
	Mail_Reg
	Mail_AddRoute
	Mail_RequestData
	Mail_ResponseData
	Mail_PushData
	Mail_ClientLeaveNotifyMe
	Mail_ClientLeaveNotification
	Mail_ClientNotFound
)

type MailNodeInfo struct {
	Address string
	Name    string
}

type MailRoutineInfo struct {
	Name string
	Routes []uint32
}

type MailClientInfo struct {
	ClientId  string
	Type      uint32 // 0 是request消息
	RequestId uint32
	Route uint32
	Data      []byte

	RemoteAddress    string
	SourceAddress    string
	SourceName    string
}

func CastMailClientInfo(obj interface{}) (MailClientInfo, error) {
	var result MailClientInfo
	switch t := obj.(type) {
	default:
		return result, errors.New("Cast To MailClientInfo Failed.")
	case MailClientInfo:
		return t, nil
	}
}