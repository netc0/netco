package frontend

import (
	"github.com/netc0/netco/connector"
	"github.com/netc0/gate/modle"
	"sync"
	"log"
)

// define


//func NewSession(inst interface{}) *ISession{
//	var ptr interface{}
//	ptr = inst
//
//}

func (this* TCPSession) IsOk() bool {
	return this.isOk
}


func SessionHandleByte(s *ISession, data []byte) {

}

var (
	allSessions map[string]*connector.Session
	tcp TCPTransporter

	sessions map[string]ISession
	sessionMutex *sync.Mutex
	DispatchBackendCallback func(s interface{}, requestId uint32, routeId uint32, data []byte)
)

type ITransporter interface {
	start()
	releaseSessions()
	checkHeartBeat()
}

type Transporter struct {
	ITransporter
	running bool   // 是否在运行中
	Host    string // 绑定的Host
	OnDataPacket func(interface{}, uint32, uint32, []byte) // 收到消息
}

type TCPTransporter struct {
	Transporter
}


func NewTransporter(inst interface{}) Transporter{
	var ptr = inst.(*Transporter)
	return *ptr
}

func StartTrasnsporter(inst interface{}) {
	var ptr = inst.(*Transporter)
	(*ptr).start()
}

// 启动前端服务
func StartFrontendSerice(config* modle.FrontendConfig) {
	allSessions = make(map[string]*connector.Session)

	sessions = make(map[string]ISession)
	sessionMutex = new(sync.Mutex)
	//context.SetTCPServerHost(config.TCPHost) // 启动 TCP 服务器
	//context.OnTCPDataCallback = OnTCPData
	//context.OnTCPNewConnection = OnTCPNewConnection

	if config.TCPHost != "" {
		var tcp TCPTransporter
		tcp.Host = config.TCPHost
		tcp.OnDataPacket = func(s interface{}, requestId uint32, routeId uint32, data []byte) {
			DispatchBackend(s, requestId, routeId, data)
		}

		go func() {
			tcp.start()
		}()
	}
}

func SetDispatchBackendCallback(callback func(session interface{}, requestId uint32, routeId uint32, data []byte)) {
	DispatchBackendCallback = callback
}

func OnTCPData(session *connector.Session, requestId uint32, routeId uint32, data []byte) {
	//BackendServiceDispatch(session, requestId, routeId, data)
}

func OnTCPNewConnection(sid string, session *connector.Session) {
	allSessions[sid] = session
}

func GetTCPSession(sid string) *connector.Session {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	var s = allSessions[sid]
	return s
}

func GetSession(sid string) ISession {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	var s = sessions[sid]
	return s
}

// 清空会话
func ClearSession() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	for k, _ := range sessions {
		delete(sessions, k)
	}
}

// 遍历会话
func ForeachSession(callback func(session ISession)) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	for _, v := range sessions {
		callback(v)
	}
}
// 添加会话
func AddSession(inst interface{}) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	session, ok := inst.(ISession)
	if ok {
		sessions[session.GetId()] = session
		log.Println("cast to ISession ok")
		return
	}

	log.Println("cast to ISession error")
}
// 删除会话
func RemoveSession(inst interface{}) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	log.Println("关闭会话", inst)
	session, ok := inst.(ISession)
	if ok {
		delete(sessions, session.GetId())
	}
}

// 传递消息到后端
func DispatchBackend(s interface{}, requestId uint32, routeId uint32, data []byte) {
	log.Println("收到数据:",requestId, routeId, string(data))
	DispatchBackendCallback(s, requestId, routeId, data)
}