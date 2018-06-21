package backend

import (
	"github.com/netc0/netco/app"
	"github.com/netc0/gate/rpc"
	"github.com/netc0/gate/modle"
	"github.com/netc0/netco/nrpc"
	"log"
	"github.com/netc0/netco/message"
	"github.com/netc0/gate/services/frontend"
)

var (
	proxy *rpc.GateProxy
)

func StartBackendService(app *app.App, config *modle.BackendConfig,
	getSessionCallback func (string)(interface{})) {
	proxy = rpc.NewGateProxy(getSessionCallback)
	proxy.AuthCode = config.Auth
	app.SetRPCServerHost(config.Host, proxy)   // 启动 RPC 服务器
}

func BackendServiceDispatch(s interface{}, requestId uint32, routeId uint32, data []byte) {
	session, ok := s.(frontend.ISession)
	if !ok {
		return
	}

	msg := nrpc.RPCGateRequest{}
	msg.RequestId = requestId
	msg.RouteId  = routeId
	msg.Data = data

	msg.ClientId = session.GetId()

	err := rpc.DispatchRequest(proxy, msg)
	if err != nil {
		log.Println(err)
		r := message.BuildSimpleMessage(404, err.Error())
		session.Response(requestId, r)
	}
}