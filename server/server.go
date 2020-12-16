package server

import (
	cpb "github.com/SailGame/Core/pb/core"
	d "github.com/SailGame/Core/data"
)

// CoreServer is derived from Grpc GameCoreServer Interface
// the required methods are implemented in separated files like room.go, account.go
type CoreServer struct
{
	cpb.UnimplementedGameCoreServer
	mStorage d.Storage 
}

// CoreServerConfig contains necessary parameters when building core server
type CoreServerConfig struct
{

}

// NewCoreServer builds a core server
func NewCoreServer(config *CoreServerConfig) (*CoreServer, error) {
	return &CoreServer{}, nil
}

