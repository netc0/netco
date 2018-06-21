package modle

import (
	"github.com/netc0/netco/app"
	"github.com/netc0/netco/nrpc"
	"net/rpc"
	"hash/crc32"
	"log"
	"errors"
	"github.com/netc0/netco/connector"
)

type GateProxy struct {
	RouteIds map[uint32]string
	RouteNames map[string]*GateRPCRecord
	Context *app.App
}

func (this *GateProxy) Init () {
	this.RouteIds = make(map[uint32]string)
	this.RouteNames = make(map[string]*GateRPCRecord)
}

func (this* GateProxy) AddRoute(info nrpc.RPCAddRouteInfo, reply *int) error {
	if info.Auth != "1111" {
		return errors.New("Auth code error")
	}
	*reply = 1
	// 注册 RPC
	delete(this.RouteNames, info.RCPRemote)
	var item = &GateRPCRecord{}
	this.RouteNames[info.RCPRemote] = item
	// 尝试连接 RPC
	client, err := rpc.Dial("tcp", info.RCPRemote)
	if err != nil {
		return err
	}
	item.client = client


	for _, r := range info.Routes {
		var crc = crc32.ChecksumIEEE([]byte(r))
		this.RouteNames[r] = item
		this.RouteIds[crc] = r
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
	return this.RouteNames[routeName]
}

func (this *GateProxy) getRouteName(routeId uint32) string {
	return this.RouteIds[routeId]
}

// 分发消息到后端
func (this *GateProxy) DispatchRequest(session *connector.Session, requestId uint32, routeId uint32, data []byte) {
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
	return this.Context.GetTCPSession(id)
}