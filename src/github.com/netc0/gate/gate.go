package main

import (
	"github.com/netc0/netco/events"
	"log"
	"flag"
	"os"
	"github.com/netc0/netco/app"
	"github.com/netc0/netco/connector"
)

const HELLO_WORLD = "helloworld"

type AppArgs struct  {
	help bool
}

var (
	appArgs AppArgs
)

func parseArgs() {
	flag.BoolVar(&appArgs.help, "h", false, "显示帮助")
	flag.Parse()
}

func processArgs() {
	if appArgs.help {
		flag.Usage()
		os.Exit(0)
	}
}

// 发送消息
type HelloWorld struct {
	Name string
}

type GoodBye struct {
	Name string
}

var instance *GateProxy

func startApp () {
	var gateApp = app.NewApp()
	instance = new(GateProxy)
	instance.init()
	instance.app = &gateApp
	gateApp.SetTCPServerHost(":9000", TCPHandler) // 启动 TCP 服务器
	gateApp.SetRPCServerHost(":9001", instance) // 启动 RPC 服务器

	gateApp.OnTCPDataCallback = OnTCPData

	listener := events.NewEventListener(myHandler)
	gateApp.EventDispatcher.AddEventListener(HELLO_WORLD, listener)

	gateApp.EventDispatcher.DispatchEvent(events.NewEvent(HELLO_WORLD, HelloWorld{Name:"Hello"}))
	gateApp.EventDispatcher.DispatchEvent(events.NewEvent(HELLO_WORLD, GoodBye{Name:"Bye"}))

	gateApp.Start()
}

func main() {
	parseArgs()   // 解析参数
	processArgs() // 处理参数
	startApp()    // 启动
}

func TCPHandler(event events.Event) {

}

func myHandler(event events.Event) {
	var hw, ok = event.Object.(HelloWorld)
	if ok {
		log.Println(event.Type, hw, event.Target)
	} else {
		log.Println("event message error", event.Type, event.Object, event.Target)
	}
}

func OnTCPData(session *connector.Session, requestId uint32, routeId uint32, data []byte) {
	instance.dispatchRequest(session, requestId, routeId, data)
}