package main

import (
	"flag"
	"os"
	"github.com/netc0/netco/app"
	"github.com/netc0/gate/modle"
	"github.com/netc0/gate/rpc"
	"github.com/netc0/gate/services/frontend"
	"github.com/netc0/gate/services/backend"
)

type AppArgs struct  {
	help    bool
	RPCAuth string
	RPCHost string

	TCPHost string
	UDPHost string
}

var (
	appArgs AppArgs
	proxy *rpc.GateProxy
)

func parseArgs() {
	flag.BoolVar(&appArgs.help, "h", false, "显示帮助")
	flag.StringVar(&appArgs.RPCAuth, "k", "netc0", "RPC 验证码")
	flag.StringVar(&appArgs.RPCHost, "r", ":9002", "RPC Host")

	flag.StringVar(&appArgs.TCPHost, "t", ":9000", "TCP Host")
	flag.StringVar(&appArgs.UDPHost, "u", ":9001", "TCP Host")
	flag.Parse()
}

func processArgs() {
	if appArgs.help {
		flag.Usage()
		os.Exit(0)
	}
}

func setupFrontend(config* modle.FrontendConfig) {
	config.TCPHost = ":9000"
	config.UDPHost = ":9001"
}
func setupBackend(config* modle.BackendConfig) {
	config.Host = appArgs.RPCHost
	config.Auth = appArgs.RPCAuth
}

func startApp () {
	var context = app.NewApp()

	// 前端配置参数
	var frontendConfig modle.FrontendConfig
	setupFrontend(&frontendConfig)

	// 后端配置
	var backendConfig modle.BackendConfig
	setupBackend(&backendConfig)

	//services.StartEventService(instance) // 事件服务
	frontend.StartFrontendSerice(&frontendConfig) // 前端服务
	frontend.SetDispatchBackendCallback(DispatchBackend)

	backend.StartBackendService(&context, &backendConfig, getSessionCallback) // 后端服务
	context.Start()
}

func DispatchBackend(s interface{}, requestId uint32, routeId uint32, data []byte) {
	backend.BackendServiceDispatch(s, requestId, routeId, data);
}

func getSessionCallback(sid string) interface{}{
	return frontend.GetSession(sid)
}

func main() {
	parseArgs()   // 解析参数
	processArgs() // 处理参数
	startApp()    // 启动
}
