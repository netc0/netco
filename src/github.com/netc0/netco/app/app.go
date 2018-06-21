package app

import (
	"github.com/netc0/netco/events"
	"time"
	"github.com/netc0/netco/connector"
	"log"
	"github.com/netc0/netco/nrpc"
)

type App struct {
	EventDispatcher *events.EventDispatcher
	aTCPTransporter *connector.TCPTransporter
	aTCPHost        string
	aTCPHander 		events.EventHandler
	aRPCHost        string
	aRPCHandler     interface{}

	OnTCPNewConnection func(id string, session *connector.Session)
	OnTCPCloseConnection func(id string)
	OnTCPDataCallback func(session *connector.Session, requestId uint32, routeId uint32, data []byte)
}

func NewApp() App {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var app App
	app.EventDispatcher = events.NewEventDispatcher()
	app.aTCPTransporter = nil

	return app
}

func (this *App) SetTCPServerHost(TCPHost string, handler events.EventHandler) {
	this.aTCPHost = TCPHost
}

func (this *App) startTCPServer() {
	this.aTCPTransporter = connector.CreateTCPConnector(this.aTCPHost)
	this.aTCPTransporter.Start(
		func(id string, session *connector.Session) {
			this.onTCPNewConnection(id, session)
		},
		func(s string) {
			this.onTCPCloseConnection(s)
		},
		func(session *connector.Session,
			requestId uint32, routeId uint32, data []byte) {
			this.onTCPData(session, requestId, routeId, data)
		})
}

func (this *App) startRPCServer() {
	if this.aRPCHost == "" {return}

	nrpc.RPCServerStart(this.aRPCHost, this.aRPCHandler)
}

func (this *App) SetRPCServerHost(RPCHost string, RPCHandler interface{}) {
	this.aRPCHost = RPCHost
	this.aRPCHandler = RPCHandler
}

func (this *App) Start () {
	go this.startTCPServer()
	go this.startRPCServer()
	log.Println("App启动完成")
	for {
		time.Sleep(time.Second)
	}
}

// tcp 新连接进入
func (this *App) onTCPNewConnection(id string, session *connector.Session) {
	if this.OnTCPNewConnection != nil {
		this.OnTCPNewConnection(id, session)
	}
}

// tcp 关闭连接
func (this *App) onTCPCloseConnection(id string) {
	if this.OnTCPCloseConnection != nil {
		this.OnTCPCloseConnection(id)
	}
}

// tcp 新数据进入, data为transporter完整的packet数据
func (this *App)onTCPData(session *connector.Session, requestId uint32, routeId uint32, data []byte) {
	log.Println("连接数据", "RequestId:", requestId, " routeId:", routeId, " Data:" , string(data))
	this.OnTCPDataCallback(session, requestId, routeId, data)
}

func (this *App) GetTCPSession(id string) (*connector.Session) {
	return this.aTCPTransporter.GetSession(id)
}