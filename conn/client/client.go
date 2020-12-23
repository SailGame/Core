package client

import (
	"sync"
	"sync/atomic"

	cpb "github.com/SailGame/Core/pb/core"
)

type Conn struct{
	mServer cpb.GameCore_ListenServer
	mWg sync.WaitGroup
	mRunning atomic.Value
}

func NewConn(server cpb.GameCore_ListenServer) (*Conn) {
	conn := &Conn{
		mServer: server,
	}
	conn.mRunning.Store(false)
	return conn
}

func (conn *Conn) Send(msg *cpb.BroadcastMsg) (error) {
	return conn.mServer.Send(msg)
}

func (conn *Conn) Start() {
	if(!conn.mRunning.Load().(bool)){
		// return err?
		conn.mWg.Add(1)
	}
	conn.mRunning.Store(true)
	conn.mWg.Wait()
}

func (conn *Conn) Close() {
	if(!conn.mRunning.Load().(bool)){
		return
	}
	conn.mRunning.Store(false)
	conn.mWg.Done()
}