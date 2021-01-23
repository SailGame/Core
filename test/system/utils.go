package system

import (
	"context"

	cpb "github.com/SailGame/Core/pb/core"
)

func buildOneUserAndOneRoom(f *systemFixture) (uc *userClient, userName string, token string, roomId int32) {
	uc = f.newUserClient()
	userName = "test"
	loginRet, _ := uc.mCoreClient.Login(context.TODO(), &cpb.LoginArgs{
		UserName: userName,
		Password: "",
	})

	token = loginRet.Token

	createRoomRet, _ := uc.mCoreClient.CreateRoom(context.TODO(), &cpb.CreateRoomArgs{
		Token: token,
	})
	roomId = createRoomRet.GetRoomId()
	uc.mCoreClient.JoinRoom(context.TODO(), &cpb.JoinRoomArgs{
		Token:  token,
		RoomId: roomId,
	})

	uc.mLisClient, _ = uc.mCoreClient.Listen(context.TODO(), &cpb.ListenArgs{Token: token})
	return
}
