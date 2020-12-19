package server

import (
	"context"

	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer CoreServer) Login(ctx context.Context, req *cpb.LoginArgs) (*cpb.LoginRet, error) {
	// TODO: User register
	coreServer.mStorage.CreateUser(req.UserName, req.Password)
	user, err := coreServer.mStorage.FindUser(req.UserName, req.Password)
	if err != nil {
		return &cpb.LoginRet{Errno: cpb.ErrorNumber_UnkownError}, nil
	}
	token, err := coreServer.mStorage.CreateToken(user)
	if err != nil {
		return &cpb.LoginRet{Errno: cpb.ErrorNumber_UnkownError}, nil
	}
	// TODO: points system
	return &cpb.LoginRet{Token: token.GetKey(), Account: &cpb.Account{UserName: user.GetUserName(), Points: 0}}, nil
}

// func (coreServer CoreServer) QueryAccount(ctx context.Context, req *cpb.QueryAccountArgs) (*cpb.QueryAccountRet, error) {
// 	return nil, nil
// }
