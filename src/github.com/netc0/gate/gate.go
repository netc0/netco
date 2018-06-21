package main

import (
	"flag"
	"os"
	"github.com/netc0/netco/app"
	"github.com/netc0/gate/services"
	"github.com/netc0/gate/modle"
	"github.com/netc0/gate/rpc"
	"github.com/netc0/netco/connector"
)

type AppArgs struct  {
	help    bool
	RPCAuth string
}

var (
	appArgs AppArgs
	instance *modle.GateProxy
	proxy *rpc.GateProxy
)

func parseArgs() {
	flag.BoolVar(&appArgs.help, "h", false, "显示帮助")
	flag.StringVar(&appArgs.RPCAuth, "k", "netc0", "RPC 验证码")
	flag.Parse()
}

func processArgs() {
	if appArgs.help {
		flag.Usage()
		os.Exit(0)
	}
}

func setupFrontend(config* modle.FrontendConfig) {
	config.Host = ":9000"
}
func setupBackend(config* modle.BackendConfig) {
	config.Host = ":9001"
	config.Auth = "12345"
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
	services.StartFrontendSerice(&context, &frontendConfig) // 前端服务
	services.StartBackendService(&context, &backendConfig, getSessionCallback) // 后端服务

	context.Start()
}

func getSessionCallback(sid string) *connector.Session{
	return services.GetTCPSession(sid)
}

func main() {
	parseArgs()   // 解析参数
	processArgs() // 处理参数
	startApp()    // 启动
}
