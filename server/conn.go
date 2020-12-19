package server

import (
	"sync"
	"sync/atomic"
)

type commonConn struct{
	mRecvLoopWg sync.WaitGroup
	mRunning atomic.Value

	// the derived conn must implement this func
	mRecvLoop func()
}

func (conn commonConn) start(){
	if(conn.mRunning.Load().(bool)){
		return
	}
	conn.mRunning.Store(true)
	conn.mRecvLoopWg.Add(1)
	go conn.mRecvLoop()
}

func (conn commonConn) stop(sync bool){
	conn.mRunning.Store(false)
	if(sync){
		conn.mRecvLoopWg.Wait()
	}
}