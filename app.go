package netco

import (
	"log"
	"os"
	"runtime"
	"os/signal"
	"syscall"
	"github.com/netc0/netco/def"
	"fmt"
	"github.com/netc0/netco/common"
)

type IApp interface {
	OnStart()
	OnDestroy()
	// events
	DispatchEvent(name string, obj interface{})
	OnEvent(name string, cb func(obj interface{})) int
	RemoveEvent(id int)
	// service
	RegisterService(name string, srv def.IService)
}

type App struct {
	eventDispatcher *EventDispatcher

	signal chan os.Signal // 信号
	Derived IApp // 实现的 app

	aRPCHost        string
	aRPCHandler     interface{}

	services map[string]def.IService
	isRunning bool
}

var (
	logger = common.GetLogger()
)

//func NewApp() App {
//	log.SetFlags(log.LstdFlags | log.Lshortfile)
//	var app App
//	app.eventDispatcher = NewEventDispatcher()
//
//	return app
//}

// 初始化
func (this *App) init() {
	this.services = make(map[string]def.IService)
	this.isRunning = false
	logger.Prefix("[app] ")
}

// 开始启动
func (this *App) Run() {
	this.init()
	runtime.GC()
	// 监听消息
	this.signal = make(chan os.Signal)
	signal.Notify(this.signal, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE)
	// 启动Event System
	this.eventDispatcher = NewEventDispatcher()

	this.isRunning = true

	this.Derived.OnStart() // 回调 OnStart

	// 等待退出消息
	for this.isRunning {
		select {
		case sig := <-this.signal:
			switch sig {
			case syscall.SIGPIPE:
			default:
				this.isRunning = false
			}
			logger.Debug("receive signal. sig=", sig)
		}
	}
	// OnDestroy
	for _, v := range this.services  {
		v.OnDestroy()
	}

	this.Derived.OnDestroy() // 回调 OnDestroy
}
//
//func (this *App) startRPCServer() {
//	if this.aRPCHost == "" {return}
//
//	RPCServerStart(this.aRPCHost, this.aRPCHandler)
//}
//
//func (this *App) SetRPCServerHost(RPCHost string, RPCHandler interface{}) {
//	this.aRPCHost = RPCHost
//	this.aRPCHandler = RPCHandler
//}
//
//func (this *App) Start () {
//	go this.startRPCServer()
//	log.Println("App启动完成")
//}

// 发送消息
func (this *App) DispatchEvent(name string, obj interface{}) {
	this.eventDispatcher.DispatchEvent(NewEvent(name, obj))
}
// 处理消息
func (this *App) OnEvent(name string, cb func(obj interface{})) int {
	l := NewEventListener( func(e Event) {
		cb(e.Object)
	})
	return this.eventDispatcher.AddEventListener(name, l)
}

// 移除消息处理
func (this *App) RemoveEvent(id int) {
	this.eventDispatcher.RemoveEventListener(id)
}

// 注册服务
func (this *App) RegisterService(name string, srv def.IService) {
	if this.services[name] != nil {
		log.Panic(fmt.Sprintf("service[%v] exist!", name))
	}
	this.services[name] = srv
	if this.isRunning {
		srv.OnStart() // 已经启动 直接运行
	}
}