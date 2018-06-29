package main

import (
	"log"
	"github.com/netc0/netco/nrpc"
)

type Example struct {
	nrpc.RPCReceiver
	context *ExampleContext
}

func (this *Example) Test(info nrpc.RPCGateRequest, reply *int) error{
	log.Println("call test")
	var response nrpc.RPCGateResponse
	response.RequestId = info.RequestId
	response.ClientId = info.ClientId
	response.Data = []byte("你好呀111")

	var rreply int
	this.context.gateRPC.Call("GateProxy.Reply", response, &rreply)

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
	this.context.gateRPC.Call("GateProxy.Reply", response, &rreply)

	*reply = 100

	this.registerGateSessionClose(response.ClientId)
	return nil
}

// 关闭连接
func (this *Example) OnSessionClose(info nrpc.RPCMessage, reply *int) error{
	log.Println("OnSessionClose", string(info.Value))
	return nil
}


func (this* Example) registerGateSessionClose(id string) {
	var msg nrpc.RPCMessage
	msg.Command = 1
	msg.Value = []byte(id)
	msg.AuthCode = this.context.auth
	msg.ResponseNodeName = this.context.nodeName
	msg.ResponseRoute = "Example.OnSessionClose"

	var r int
	c := this.context.gateRPC.Call("GateProxy.OnMessage", msg, &r)
	if c != nil {
		log.Println(c.Error())
	}
}

func (this *Example) OnMessage(message nrpc.RPCMessage) {
	log.Println("OnRPCMessage", message.Command)
}