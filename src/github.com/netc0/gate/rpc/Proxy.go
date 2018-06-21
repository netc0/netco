package rpc

import (
	"github.com/netc0/netco/nrpc"
	"errors"
	"net/rpc"
	"log"
	"sync"
	"time"
	"fmt"
	"github.com/netc0/netco/connector"
	"hash/crc32"
)

type BackendInfo struct{
	nodeName     string
	client       * rpc.Client
	heatBeatTime time.Time
}

type GateProxy struct {
	AuthCode string
	Backends map[string]*BackendInfo
	backendLock *sync.Mutex
	getSessionCallback func (string)(*connector.Session)
	backendCache map[uint32]string // crc -> backend name
	backendRoute map[uint32]string // crc -> route
}

//
func NewGateProxy(getSessionCallback func (string)(*connector.Session)) *GateProxy {
	var v GateProxy
	v.init()
	v.getSessionCallback = getSessionCallback
	return &v
}

// 初始化
func (this *GateProxy) init() {
	this.Backends = make(map[string]*BackendInfo)
	this.backendLock = new(sync.Mutex)
	this.backendCache = make(map[uint32]string)
	this.backendRoute = make(map[uint32]string)
	//
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for range ticker.C {
			go this.checkHeartBeat()
		}
	}()
}

// RPC 注册后端
func (this *GateProxy) RegisterBackend(info nrpc.RPCBackendInfo, reply* int) error {
	log.Println("RegisterBackend", info)
	if this.AuthCode != info.Auth {
		return errors.New("Auth Code Invalid.")
	}
	return this.addBackend(info)
}
// RPC 回复
func (this *GateProxy) Reply(info nrpc.RPCGateResponse, reply* int) error {
	return this.reply(info)
}

// BackendHeartBeat
func (this *GateProxy) BackendHeartBeat(info nrpc.RPCBackendInfo, reply* int) error {
	if this.AuthCode != info.Auth {
		return errors.New("Auth Code Invalid.")
	}

	this.backendLock.Lock()
	defer this.backendLock.Unlock()

	var backend = this.getBackend(info.Name)
	if backend == nil {
		return errors.New("Backend not register")
	}
	backend.heatBeatTime = time.Now()
	return nil
}

// 注册后端
func (this *GateProxy) addBackend(info nrpc.RPCBackendInfo) error {
	// 注册路由
	log.Println("注册后端", info)
	this.backendLock.Lock()
	defer this.backendLock.Unlock()
	if this.Backends[info.Name] != nil {
		delete(this.Backends, info.Name)
	}

	cli, err := nrpc.RPCClientConnect(info.RCPRemote)
	if err != nil {
		return err
	}
	var backend BackendInfo
	this.Backends[info.Name] = &backend
	backend.client = cli
	backend.heatBeatTime = time.Now()
	backend.nodeName = info.Name
	// 注册 cache
	for _, route := range info.Routes {
		var crc = crc32.ChecksumIEEE([]byte(route))
		this.backendRoute[crc] = route
		this.backendCache[crc] = info.Name
	}
	return nil
}

// 回复客户端
func (this *GateProxy) reply(info nrpc.RPCGateResponse) error {
	var session = this.getSession(info.ClientId)
	if session == nil {
		return errors.New("Client session not found.")
	}
	session.Response(connector.PacketType_DATA, info.RequestId, info.Data)
	return nil
}
func (this *GateProxy) push(info nrpc.RPCGateResponse) error {
	var session = this.getSession(info.ClientId)
	if session == nil {
		return errors.New("Client session not found.")
	}
	session.Push(connector.PacketType_PUSH, info.Data)
	return nil
}


// 获取后端
func (this *GateProxy) getBackend(name string) *BackendInfo {
	return this.Backends[name]
}
// 获取客户端
func (this *GateProxy) getSession(name string) *connector.Session {
	return this.getSessionCallback(name)
}

// 心跳检测
func (this *GateProxy) checkHeartBeat() {
	this.backendLock.Lock()
	defer this.backendLock.Unlock()
	for k, backend := range this.Backends {
		if time.Now().Second() - backend.heatBeatTime.Second() > 5 {
			log.Println("Backend heartbeat fail:", backend)
			delete(this.Backends, k)

		}
	}
}
// 根据CRC获取路由, 返回 routeName, 后端名称
func (this *GateProxy) crcBackend(crc uint32) (string, *BackendInfo) {
	var backendName = this.backendCache[crc]
	if backendName == "" {
		return "", nil
	}
	return this.backendRoute[crc], this.Backends[backendName]
}

func DispatchRequest(this *GateProxy, msg nrpc.RPCGateRequest) error {
	var routeName, backend = this.crcBackend(msg.RouteId)
	if routeName == "" || backend == nil  {
		return errors.New("Backend not exist")
	}

	rs := backend.client.Call(routeName, msg, nil)
	if rs != nil {
		return errors.New(fmt.Sprintf("Backend Invoke Error:%v", rs.Error()))
	}

	return nil
}
