package nrpc

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
	log.Println("RPCServerStart On:", host)
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
	cli, err := rpc.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return cli, nil
}


