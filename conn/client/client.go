package client


import (
	cpb "github.com/SailGame/Core/pb/core"
)

type Conn struct{
	mServer cpb.GameCore_ListenServer
}

func NewConn(server cpb.GameCore_ListenServer) (*Conn) {
	conn := &Conn{
		mServer: server,
	}
	return conn
}

func (conn Conn) Send(msg *cpb.BroadcastMsg) (error) {
	return conn.mServer.Send(msg)
}