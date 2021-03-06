package system

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/SailGame/Core/data/memory"
	"github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/server"
)

func TestRoom(t *testing.T) {
	Convey("Room", t, func() {
		f := newFixture()
		f.init(&server.CoreServerConfig{
			MStorage: memory.NewStorage(),
		})
		uc := f.newUserClient()
		userName := "test"
		loginRet, err := uc.mCoreClient.Login(context.TODO(), &core.LoginArgs{
			UserName: userName,
			Password: "",
		})
		So(err, ShouldBeNil)

		token := loginRet.Token

		createRoomRet, err := uc.mCoreClient.CreateRoom(context.TODO(), &core.CreateRoomArgs{
			Token: token,
		})

		So(err, ShouldBeNil)
		So(createRoomRet.Err, ShouldEqual, core.ErrorNumber_OK)

		listRoomRet, err := uc.mCoreClient.ListRoom(context.TODO(), &core.ListRoomArgs{
			GameName: "", // show all room
		})

		So(err, ShouldBeNil)
		So(listRoomRet.Err, ShouldEqual, core.ErrorNumber_OK)
		So(len(listRoomRet.GetRoom()), ShouldEqual, 1)

		room := listRoomRet.GetRoom()[0]

		So(len(room.UserName), ShouldEqual, 0)

		joinRoomRet, err := uc.mCoreClient.JoinRoom(context.TODO(), &core.JoinRoomArgs{
			Token:  token,
			RoomId: room.GetRoomId(),
		})

		So(err, ShouldBeNil)
		So(joinRoomRet.Err, ShouldEqual, core.ErrorNumber_OK)

		qryAccountRet, err := uc.mCoreClient.QueryAccount(context.TODO(), &core.QueryAccountArgs{
			Key: &core.QueryAccountArgs_Token{
				Token: token,
			},
		})

		So(err, ShouldBeNil)
		So(qryAccountRet.Err, ShouldEqual, core.ErrorNumber_OK)
		So(qryAccountRet.GetAccount().UserName, ShouldEqual, userName)
		So(qryAccountRet.RoomId, ShouldEqual, room.GetRoomId())

		listRoomRet2, err := uc.mCoreClient.ListRoom(context.TODO(), &core.ListRoomArgs{
			GameName: "", // show all room
		})

		So(err, ShouldBeNil)
		So(listRoomRet2.Err, ShouldEqual, core.ErrorNumber_OK)
		So(len(listRoomRet.GetRoom()), ShouldEqual, 1)

		room2 := listRoomRet2.GetRoom()[0]

		So(len(room2.UserName), ShouldEqual, 1)
		So(room2.UserName[0], ShouldEqual, userName)
	})
}
