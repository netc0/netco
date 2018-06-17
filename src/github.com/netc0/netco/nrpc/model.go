package nrpc

type RPCAddRouteInfo struct {
	RCPRemote  string
	Routes []string
}

type RPCGateRequest struct {
	ClientId  string
	RequestId uint32
	RouteId   uint32
	Data 	  []byte
}

type RPCGateResponse struct {
	ClientId  string
	RequestId uint32
	Data 	  []byte
}