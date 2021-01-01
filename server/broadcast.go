package server

import (
	"github.com/SailGame/Core/conn/client"
	log "github.com/sirupsen/logrus"

	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer *CoreServer) Listen(req *cpb.ListenArgs, lisServer cpb.GameCore_ListenServer) error {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return err
	}
	log.Infof("User(%s) Listen: %v", token.GetUserName(), req)
	conn := client.NewConn(lisServer)
	token.GetUser().SetConn(conn)
	conn.Start()
	return nil
}
