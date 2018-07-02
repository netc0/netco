package main

import (
	"log"
	"github.com/netc0/netco/nrpc"
	"fmt"
	"time"
	"github.com/netc0/netco/connector"
)

type Example struct {
	context *ExampleContext
}

type Client struct {
	id string
	online bool
}

var (
	isInit bool
	clients map[string] *Client
)

func getClient(sid string, auto bool) *Client{
	if !isInit {
		clients = make(map[string] *Client)
		isInit = true
	}
	if val, ok := clients[sid]; ok {
		return val
	}
	if auto {
		cli := Client{id: sid, online: false}
		clients[sid] = &cli
		return &cli
	}
	return nil
}

func isOnline (sid string) bool {
	if cli := getClient(sid, false); cli != nil {
		return cli.online
	}

	return false
}

func (this *Example) Test(info nrpc.RPCGateRequest, reply *int) error{
	log.Println("call test", string(info.Data))
	var response nrpc.RPCGateResponse
	response.RequestId = info.RequestId
	response.ClientId = info.ClientId
	response.Data = []byte("你好呀111")

	var rreply int
	this.context.gateRPC.Call("GateProxy.Reply", response, &rreply)

	go func() {
		ticker := time.NewTicker(time.Second * 3)
		count := 0
		for range ticker.C {
			// 每3秒推送一条消息
			if isOnline(info.ClientId) {
				var data nrpc.RPCGatePush
				data.ClientId = info.ClientId
				data.Data = connector.PacketPushToBinary("game.OnPush", []byte(fmt.Sprintf("推送消息 %v", count)))

				var r int
				if this.context.gateRPC != nil {
					this.context.gateRPC.Call("GateProxy.Push", data, &r)
					count++;
				}
			} else {
				log.Println(info.ClientId, "已经离线")
				ticker.Stop()
			}
		}
	}()

	*reply = 100
	return nil
}
var msgId int = 100
func (this *Example) Login(info nrpc.RPCGateRequest, reply *int) error{
	log.Println("call login", string(info.Data))
	var response nrpc.RPCGateResponse
	response.RequestId = info.RequestId
	response.ClientId = info.ClientId
	response.Data = []byte(fmt.Sprintf("%v什么哦2222", msgId))
	msgId++;

	var rreply int
	this.context.gateRPC.Call("GateProxy.Reply", response, &rreply)

	*reply = 100

	this.registerGateSessionClose(response.ClientId)

	cli := getClient(info.ClientId, true)
	cli.online = true

	p := getClient(info.ClientId, false)
	if p != nil {
		log.Println("online?", p.online)
	}

	return nil
}

// 关闭连接
func (this *Example) OnSessionClose(info nrpc.RPCMessage, reply *int) error{
	log.Println("OnSessionClose", string(info.Value))
	sid := string(info.Value)
	cli := getClient(sid, false)
	if cli != nil {
		cli.online = false
	}
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
