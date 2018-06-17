package main

import (
	"github.com/netc0/netco/nrpc"
	"log"
	"net/rpc"
	"github.com/netc0/netco/connector"
	"hash/crc32"
	"github.com/netc0/netco/app"
	"errors"
)

type GateRPCRecord struct {
	remote string
	client *rpc.Client
	routes []string
}

type GateProxy struct {
	routeIds map[uint32]string
	routeNames map[string]*GateRPCRecord
	app *app.App
}

func (this *GateProxy) init () {
	this.routeIds = make(map[uint32]string)
	this.routeNames = make(map[string]*GateRPCRecord)
}

func (this* GateProxy) AddRoute(info nrpc.RPCAddRouteInfo, reply *int) error {
	*reply = 1
	// 注册 RPC
	delete(this.routeNames, info.RCPRemote)
	var item = &GateRPCRecord{}
	this.routeNames[info.RCPRemote] = item
	// 尝试连接 RPC
	client, err := rpc.Dial("tcp", info.RCPRemote)
	if err != nil {
		return err
	}
	item.client = client


	for _, r := range info.Routes {
		var crc = crc32.ChecksumIEEE([]byte(r))
		this.routeNames[r] = item
		this.routeIds[crc] = r
	}
	item.routes = append(item.routes, info.Routes...)
	log.Println("注册 RPC Route", info)
	return nil
}

func (this* GateProxy) Reply(info nrpc.RPCGateResponse, reply *int) error {
	// 回复客户端
	log.Println("回复 RPC Route", info, string(info.Data))
	var session = this.getSession(info.ClientId)
	if session == nil {
		log.Println("客户端已经断开")
		return errors.New("客户端已经断开")
	}

	session.Response(connector.PacketType_DATA, info.RequestId, info.Data)
	return nil
}

func (this *GateProxy) getBackend(routeName string)* GateRPCRecord{
	return this.routeNames[routeName]
}

func (this *GateProxy) getRouteName(routeId uint32) string {
	return this.routeIds[routeId]
}

// 分发消息到后端
func (this *GateProxy) dispatchRequest(session *connector.Session, requestId uint32, routeId uint32, data []byte) {
	msg := nrpc.RPCGateRequest{}
	msg.RequestId = requestId
	msg.RouteId  = routeId
	msg.Data = data
	msg.ClientId = session.GetId()

	var routeName = this.getRouteName(routeId)
	var backend = this.getBackend(routeName)
	if backend == nil {
		log.Println("后端不存在")
		return
	}
	var reply int

	var rs = backend.client.Call(routeName, msg, &reply)
	log.Println("reply", reply, rs)
}

func (this *GateProxy) getSession(id string) (*connector.Session){
	return this.app.GetTCPSession(id)
}