package main

import (
	"log"
	"github.com/netc0/netco/nrpc"
)

type Example struct {
	context *ExampleContext
}

func (this *Example) Test(info nrpc.RPCGateRequest, reply *int) error{
	log.Println("call test")
	var response nrpc.RPCGateResponse
	response.RequestId = info.RequestId
	response.ClientId = info.ClientId
	response.Data = []byte("你好呀111")

	var rreply int
	var rs = this.context.gateRPC.Call("GateProxy.Reply", response, &rreply)
	log.Println(rs)
	*reply = 100
	return nil
}

func (this *Example) Login(info nrpc.RPCGateRequest, reply *int) error{
	log.Println("call login")
	var response nrpc.RPCGateResponse
	response.RequestId = info.RequestId
	response.ClientId = info.ClientId
	response.Data = []byte("什么哦2222")

	var rreply int
	var rs = this.context.gateRPC.Call("GateProxy.Reply", response, &rreply)
	log.Println(rs)
	*reply = 100
	return nil
}
