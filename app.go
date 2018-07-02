package netco

import (
	"time"
	"log"
)

type App struct {
	EventDispatcher *EventDispatcher

	aRPCHost        string
	aRPCHandler     interface{}

}

func NewApp() App {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var app App
	app.EventDispatcher = NewEventDispatcher()

	return app
}

func (this *App) startRPCServer() {
	if this.aRPCHost == "" {return}

	RPCServerStart(this.aRPCHost, this.aRPCHandler)
}

func (this *App) SetRPCServerHost(RPCHost string, RPCHandler interface{}) {
	this.aRPCHost = RPCHost
	this.aRPCHandler = RPCHandler
}

func (this *App) Start () {
	go this.startRPCServer()
	log.Println("App启动完成")
	for {
		time.Sleep(time.Second)
	}
}
