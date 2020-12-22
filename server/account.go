package server

import (
	"context"

	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer *CoreServer) Login(ctx context.Context, req *cpb.LoginArgs) (*cpb.LoginRet, error) {
	// TODO: User register
	err := coreServer.mStorage.CreateUser(req.UserName, req.Password)
	if err != nil {
		return &cpb.LoginRet{Errno: cpb.ErrorNumber_User_FailToCreateUser}, nil
	}
	user, err := coreServer.mStorage.FindUser(req.UserName, req.Password)
	if err != nil {
		return &cpb.LoginRet{Errno: cpb.ErrorNumber_User_FailToFindUser}, nil
	}
	token, err := coreServer.mStorage.CreateToken(user)
	if err != nil {
		return &cpb.LoginRet{Errno: cpb.ErrorNumber_User_FailToGenerateToken}, nil
	}
	// TODO: rank system
	return &cpb.LoginRet{Token: token.GetKey(), Account: &cpb.Account{UserName: user.GetUserName(), Points: 0}}, nil
}

// func (coreServer *CoreServer) QueryAccount(ctx context.Context, req *cpb.QueryAccountArgs) (*cpb.QueryAccountRet, error) {
// 	return nil, nil
// }
