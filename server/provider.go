package server

import (
	"github.com/SailGame/Core/conn/provider"
	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer CoreServer) Provider(pServer cpb.GameCore_ProviderServer) error {
	conn := provider.NewConn(pServer, coreServer)
	conn.Start()
	return nil
}

func (coreServer CoreServer) HandleRegisterArgs(providerMsg *cpb.ProviderMsg, regArgs *cpb.RegisterArgs) error {
	// coreServer.mStorage.RegisterProvider(regArgs.Id)
	return nil
}

func (coreServer CoreServer) HandleNotifyMsg(*cpb.ProviderMsg, *cpb.NotifyMsg) error {
	return nil
}

func (coreServer CoreServer) HandleStartGameRet(*cpb.ProviderMsg, *cpb.StartGameRet) error {
	return nil
}

func (coreServer CoreServer) HandleQueryStateRet(*cpb.ProviderMsg, *cpb.QueryStateRet) error {
	return nil
}

func (coreServer CoreServer) Disconnect(conn *provider.Conn) {
	coreServer.mStorage.UnRegisterProvider(conn.ID)
	if(conn.ID != ""){
		coreServer.mProviders.Delete(conn.ID)
	}
}
