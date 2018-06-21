package modle

import "net/rpc"

type GateRPCRecord struct {
	remote string
	client *rpc.Client
	routes []string
}
