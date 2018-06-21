package services

import (
	"github.com/netc0/netco/app"
	"github.com/netc0/netco/events"
	"log"
	"github.com/netc0/netco/connector"
	"github.com/netc0/gate/modle"
)

var (
	allSessions map[string]*connector.Session
)


// 启动前端服务
func StartFrontendSerice(context *app.App, config* modle.FrontendConfig) {
	allSessions = make(map[string]*connector.Session)

	context.SetTCPServerHost(config.Host, TCPHandler) // 启动 TCP 服务器
	context.OnTCPDataCallback = OnTCPData
	context.OnTCPNewConnection = OnTCPNewConnection
}

func TCPHandler(event events.Event) {
	log.Println(event)
}

func OnTCPData(session *connector.Session, requestId uint32, routeId uint32, data []byte) {
	BackendServiceDispatch(session, requestId, routeId, data)
}

func OnTCPNewConnection(sid string, session *connector.Session) {
	allSessions[sid] = session
}

func GetTCPSession(sid string) *connector.Session {
	return allSessions[sid]
}