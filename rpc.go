package netco

import (
	"net"
	"log"
	"net/rpc"
)

func RPCServerStart(host string, rscv interface{}) error {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		log.Println(err)
		return err
	}
	rpc.Register(rscv)
	log.Println("Frontend启动 RPC:", host)
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(c)
		}
	}()
	return nil
}

func RPCClientConnect(host string) (*rpc.Client, error){
	cli, err := rpc.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	return cli, nil
}


