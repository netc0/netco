package main

import (
	"log"
	"net/rpc"
	"time"
	"github.com/netc0/netco/nrpc"
	"sync"
)

type ExampleContext struct {
	gateRPC *rpc.Client
	gateLock *sync.Mutex

	nodeName string
	auth string
}

func runRPCServer(context *ExampleContext) {
	var v = new(Example)
	v.context = context;
	nrpc.RPCServerStart(":8001", v)
}

func connectGate(context *ExampleContext) {
	context.gateLock.Lock()
	defer context.gateLock.Unlock()

	if context.gateRPC != nil {
		// 已经连接 connected
		return
	}

	cli, err := nrpc.RPCClientConnect("127.0.0.1:9002")
	if err != nil {
		log.Println(err)
		return
	}
	context.gateRPC = cli
	log.Println("注册 Proxy")

	var info nrpc.RPCBackendInfo
	info.RCPRemote = "127.0.0.1:8001"

	info.Name = context.nodeName
	info.Auth = context.auth
	info.Routes = append(info.Routes, "Example.Test")
	info.Routes = append(info.Routes, "Example.Login")

	reply := 0
	rs := context.gateRPC.Call("GateProxy.RegisterBackend", info, &reply)
	log.Println("GateProxy reply:", rs)
}

// 监控网关的连接状态
func gateMonitor(context *ExampleContext) {
	ticker := time.NewTicker(time.Second * 3)
	for range ticker.C {
		if context.gateRPC == nil {
			go connectGate(context)
		} else {
			go gateHeartBeat(context)
		}
	}
}

func gateHeartBeat(context *ExampleContext) {
	if context.gateRPC != nil {
		var info nrpc.RPCBackendInfo
		info.Name = context.nodeName
		info.Auth = context.auth
		rs := context.gateRPC.Call("GateProxy.BackendHeartBeat", info, nil)
		if rs != nil { // 断开连接
			log.Println(rs.Error())
			context.gateRPC = nil
		}
	}
}

func main() {
	var ctx ExampleContext
	ctx.gateLock = new(sync.Mutex)
	ctx.nodeName = "exampleGame"
	ctx.auth = "netc0"
	log.Println("启动游戏")

	go runRPCServer(&ctx)
	go connectGate(&ctx)
	go gateMonitor(&ctx)

	for {
		time.Sleep(time.Second)
	}
}
