package provider

import (
	"log"
	"sync"
	"sync/atomic"

	cpb "github.com/SailGame/Core/pb/core"
)

type Handler interface {
	HandleRegisterArgs(*Conn, *cpb.ProviderMsg, *cpb.RegisterArgs) error
	HandleNotifyMsg(*Conn, *cpb.ProviderMsg, *cpb.NotifyMsg) error

	Disconnect(*Conn)
}

type Conn struct {
	// bind to an entity
	ID         interface{}
	// used in logging
	PrintID    string

	mRecvLoopWg sync.WaitGroup
	mRunning    atomic.Value
	mServer     cpb.GameCore_ProviderServer
	mHandler    Handler
}

func NewConn(pServer cpb.GameCore_ProviderServer, handler Handler) *Conn {
	conn := &Conn{
		mServer: pServer,
		mHandler: handler,
	}
	conn.mRunning.Store(false)
	return conn
}

func (conn Conn) RecvLoop() {
	defer conn.mRecvLoopWg.Done()
	for conn.mRunning.Load().(bool) {
		msg, err := conn.mServer.Recv()
		if err != nil {
			conn.mHandler.Disconnect(&conn)
			conn.mRunning.Store(false)
			return
		}

		if submsg := msg.GetRegisterArgs(); submsg != nil{
			conn.mHandler.HandleRegisterArgs(&conn, msg, submsg)
		}else if submsg := msg.GetNotifyMsg(); submsg != nil{
			conn.mHandler.HandleNotifyMsg(&conn, msg, submsg)
		}else{
			log.Printf("Received unwanted msg (%s) from provider %s", msg.String(), conn.PrintID)
			conn.mHandler.Disconnect(&conn)
			conn.mRunning.Store(false)
			return
		}
	}
}

func (conn Conn) Start() {
	if conn.mRunning.Load().(bool) {
		return
	}
	conn.mRunning.Store(true)
	conn.mRecvLoopWg.Add(1)
	go conn.RecvLoop()
}

func (conn Conn) Stop(sync bool) {
	conn.mRunning.Store(false)
	if sync {
		conn.mRecvLoopWg.Wait()
	}
}

func (conn Conn) Send(msg *cpb.ProviderMsg) error {
	return conn.mServer.Send(msg)
}