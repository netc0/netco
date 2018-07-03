package netco

import (
	"net/rpc"
	"log"
	"errors"
	"fmt"
)

type backendInfo struct {
	gateRPC *rpc.Client
}

func GateConnect(gateHost string, gateAuth string, node string, backend string, auth string, routes []string) (interface{}, error) {
	var info RPCBackendInfo
	var result interface{}

	info.RCPRemote = backend

	info.AuthCode = gateAuth

	info.Name = node
	info.Auth = auth
	info.Routes = routes

	if node == "" {
		return result, errors.New("Node Name Is Nil")
	}

	var be backendInfo

	log.Println("connect:", info.AuthCode)
	cli, err := RPCClientConnect(gateHost)
	if err != nil {
		return result, err
	}
	be.gateRPC = cli

	reply := 0
	rs := cli.Call("GateProxy.RegisterBackend", info, &reply)
	if rs != nil {
		log.Println("GateProxy reply:", rs)
		return result, errors.New(fmt.Sprintf("%v", rs))
	}
	result = &be
	return result, nil
}

func getbackendInfo(itf interface{}) (*backendInfo, error){
	var be *backendInfo = nil
	switch t := itf.(type) {
	default:
		return nil, errors.New("interface args not correct")
	case *backendInfo:
		be = t
	}
	if be == nil {
		return nil, errors.New("gate lost connect")
	}
	return be, nil
}

func GateHeartBeat(itf interface{}, node string, auth string) error {
	be, err := getbackendInfo(itf)
	if err != nil {
		return err;
	}

	var info RPCBackendInfo
	info.Name = node
	info.AuthCode = auth
	rs := be.gateRPC.Call("GateProxy.BackendHeartBeat", info, nil)
	if rs != nil { // 断开连接
		be.gateRPC = nil
		return errors.New(rs.Error())
	}
	return nil
}

func GateEventDisconnect(itf interface{}, node string, auth string, response string, sid string) error {
	be, err := getbackendInfo(itf)
	if err != nil {
		return err;
	}
	var msg RPCMessage
	msg.Command = 1
	msg.Value = []byte(sid)
	msg.AuthCode = auth
	msg.ResponseNodeName = node
	msg.ResponseRoute = response

	var r int
	c := be.gateRPC.Call("GateProxy.OnMessage", msg, &r)
	if c != nil {
		return errors.New(c.Error())
	}
	return nil
}

func GatePushData(itf interface{}, sid string, route string, data []byte) error {
	be, err := getbackendInfo(itf)
	if err != nil {
		return err;
	}
	var msg RPCGatePush
	msg.ClientId = sid
	msg.Data = BuildPushPacket(route, data)

	var r int
	c := be.gateRPC.Call("GateProxy.Push", msg, &r)
	if c != nil {
		return errors.New(c.Error())
	}
	return nil
}

func GateReply(itf interface{}, sid string, requestId uint32, data []byte) error {
	be, err := getbackendInfo(itf)
	if err != nil {
		return err;
	}
	var msg RPCGatePush
	msg.ClientId = sid
	msg.Data = BuildReplyPacket(requestId, data)

	var r int
	c := be.gateRPC.Call("GateProxy.Push", msg, &r)
	if c != nil {
		return errors.New(c.Error())
	}
	return nil
}