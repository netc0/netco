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

	// base info
	GetNodeName() string
	GetNodeAddress() string
	GetGateAddress() string
	SetNodeName(string)
	SetNodeAddress(string)
	SetGateAddress(string)
}

type App struct {
	eventDispatcher *EventDispatcher

	signal chan os.Signal // 信号
	Derived IApp // 实现的 app

	aRPCHost        string
	aRPCHandler     interface{}

	services map[string]def.IService
	isRunning bool

	mName string
	mAddress string
	mGateAddress string
}

var (
	logger = common.GetLogger()
)

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

// 节点名称
func (this *App) GetNodeName() string{
	return this.mName
}
func (this *App) SetNodeName(n string) {
	this.mName = n
}

// 节点地址
func (this *App) GetNodeAddress() string{
	return this.mAddress
}
func (this *App) SetNodeAddress(n string) {
	this.mAddress = n
}
// 网关地址
func (this *App) GetGateAddress() string{
	return this.mGateAddress
}
func (this *App) SetGateAddress(n string) {
	this.mGateAddress = n
}