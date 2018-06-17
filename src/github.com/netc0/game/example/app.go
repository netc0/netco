package main

import (
	"log"
	"net/rpc"
	"time"
	"github.com/netc0/netco/nrpc"
)

type ExampleContext struct {
	gateRPC *rpc.Client
}

func runRPCServer(context *ExampleContext) {
	var v = new(Example)
	v.context = context;
	nrpc.RPCServerStart(":8001", v)
}

func connectGate(context *ExampleContext) {
	cli, err := nrpc.RPCClientConnect("127.0.0.1:9001")
	if err != nil {
		log.Println(err)
		return
	}
	context.gateRPC = cli
	log.Println("测试 RPC")
	var reply int
	var info nrpc.RPCAddRouteInfo
	info.RCPRemote = "127.0.0.1:8001"

	info.Routes = append(info.Routes, "Example.Test")
	info.Routes = append(info.Routes, "Example.Login")

	log.Println(info)

	context.gateRPC.Call("GateProxy.AddRoute", info, &reply)
	log.Println("reply:", reply)
}

func main() {
	var ctx ExampleContext

	log.Println("启动游戏")

	go runRPCServer(&ctx)
	go connectGate(&ctx)

	for {
		time.Sleep(time.Second)
	}
}
