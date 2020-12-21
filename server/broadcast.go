package server

import (
	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer CoreServer) Listen(req *cpb.ListenArgs, lisServer cpb.GameCore_ListenServer) (error) {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if(err != nil){
		return err
	}
	conn := &userConn{mServer: lisServer}
	coreServer.mClients.Store(token.GetUser().GetUserName(), conn)
	return nil
}