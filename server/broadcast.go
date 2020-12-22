package server

import (
	"github.com/SailGame/Core/conn/client"

	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer *CoreServer) Listen(req *cpb.ListenArgs, lisServer cpb.GameCore_ListenServer) error {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return err
	}
	conn := client.NewConn(lisServer)
	token.GetUser().SetConn(conn)
	return nil
}
