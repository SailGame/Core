package server

import (
	"context"
	log "github.com/sirupsen/logrus"

	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer *CoreServer) Login(ctx context.Context, req *cpb.LoginArgs) (*cpb.LoginRet, error) {
	// TODO: User register
	log.Infof("Login: %v", req)
	if !coreServer.mStorage.IsUserExist(req.UserName){
		err := coreServer.mStorage.CreateUser(req.UserName, req.Password)
		if err != nil {
			return &cpb.LoginRet{Err: cpb.ErrorNumber_User_FailToCreateUser}, nil
		}
	}

	user, err := coreServer.mStorage.FindUser(req.UserName, req.Password)
	if err != nil {
		return &cpb.LoginRet{Err: cpb.ErrorNumber_User_FailToFindUser}, nil
	}
	token, err := coreServer.mStorage.CreateToken(user)
	if err != nil {
		return &cpb.LoginRet{Err: cpb.ErrorNumber_User_FailToGenerateToken}, nil
	}
	// TODO: rank system
	return &cpb.LoginRet{Token: token.GetKey(), Account: &cpb.Account{UserName: user.GetUserName(), Points: 0}}, nil
}

func (coreServer *CoreServer) QueryAccount(ctx context.Context, req *cpb.QueryAccountArgs) (*cpb.QueryAccountRet, error) {
	if req.GetUserName() != "" {
		if !coreServer.mStorage.IsUserExist(req.GetUserName()){
			return &cpb.QueryAccountRet{Err: cpb.ErrorNumber_QueryAccount_InvalidUserNameOrToken}, nil
		}
		return &cpb.QueryAccountRet{
			Err: cpb.ErrorNumber_OK,
			Account: &cpb.Account{
				UserName: "",
				Points: 0,
			},
		}, nil
	}else if req.GetToken() != "" {
		token, err := coreServer.mStorage.FindToken(req.GetToken())
		if err != nil {
			return &cpb.QueryAccountRet{Err: cpb.ErrorNumber_QueryAccount_InvalidUserNameOrToken}, nil
		}
		user := token.GetUser()
		user.Lock()
		defer user.Unlock()
		room, err := user.GetRoom()
		var roomID int32 = -1
		if err == nil {
			roomID = room.GetRoomID()
		}
		return &cpb.QueryAccountRet{
			Err: cpb.ErrorNumber_OK,
			Account: &cpb.Account{
				UserName: user.GetUserName(),
				Points: 0, }, 
				RoomId: roomID,
		}, nil
	}else {
		return &cpb.QueryAccountRet{Err: cpb.ErrorNumber_QueryAccount_InvalidUserNameOrToken}, nil
	}
}
