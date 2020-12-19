package server

import (
	cpb "github.com/SailGame/Core/pb/core"
)

type providerHandler interface{
	handle(*cpb.ProviderMsg) (error)
	disconnect(*providerConn)
}
type providerConn struct{
	commonConn
	mID string
	mServer cpb.GameCore_ProviderServer
	mHandler providerHandler
}

func newProviderConn(pServer cpb.GameCore_ProviderServer) (*providerConn){
	conn := &providerConn{
		mServer: pServer,
	}
	conn.mRunning.Store(false)
	conn.mRecvLoop = conn.recvLoop
	return conn
}

func (coreServer CoreServer) Provider(pServer cpb.GameCore_ProviderServer) error {
	conn := newProviderConn(pServer)
	conn.start()
	return nil
}

func (conn providerConn) recvLoop(){
	defer conn.mRecvLoopWg.Done()
	for conn.mRunning.Load().(bool) {
		msg, err := conn.mServer.Recv()
		if(err != nil){
			conn.mHandler.disconnect(&conn)
			conn.mRunning.Store(false)
			return
		}
		conn.mHandler.handle(msg)
	}
}