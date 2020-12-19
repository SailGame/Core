package server

import (
	"sync"

	d "github.com/SailGame/Core/data"
	cpb "github.com/SailGame/Core/pb/core"
)

type userConn struct{
	mServer cpb.GameCore_ListenServer
}

// CoreServer is derived from Grpc GameCoreServer Interface
// the required methods are implemented in separated files like room.go, account.go
type CoreServer struct
{
	cpb.UnimplementedGameCoreServer
	mStorage d.Storage
	mClients sync.Map // userName -> userConn
	mProviders sync.Map // providerID -> providerConn
}

// CoreServerConfig contains necessary parameters when building core server
type CoreServerConfig struct
{
	mStorage d.Storage
}

// NewCoreServer builds a core server
func NewCoreServer(config *CoreServerConfig) (*CoreServer, error) {
	cs := CoreServer{}
	cs.mStorage = config.mStorage
	return &cs, nil
}

