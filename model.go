package netco

type RPCAddRouteInfo struct {
	RCPRemote  string
	Routes []string
	Auth      string
}

// 客户端请求服务端消息
type RPCGateRequest struct {
	ClientId  string
	RequestId uint32
	RouteId   uint32
	Data 	  []byte
}

// 服务端回复客户端消息
type RPCGateResponse struct {
	ClientId  string
	RequestId uint32
	Data 	  []byte
}

// 推送消息
type RPCGatePush struct {
	ClientId  string
	Data 	  []byte
}

/// 后端模型
type RPCBackendInfo struct {
	Name string
	RCPRemote  string
	Routes []string
	Auth      string
}

// RPC 消息
type RPCMessage struct {
	Command  int
	AuthCode string
	Value    []byte
	ResponseNodeName string
	ResponseRoute    string
	ResponseAuthCode string
}