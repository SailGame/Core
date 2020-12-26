package system

import (
	"context"
	"testing"

	"github.com/smartystreets/assertions"
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
		So(err, assertions.ShouldBeNil)

		token := loginRet.Token

		createRoomRet, err := uc.mCoreClient.CreateRoom(context.TODO(), &core.CreateRoomArgs{
			Token: token,
		})

		So(err, assertions.ShouldBeNil)
		So(createRoomRet.Err, assertions.ShouldEqual, core.ErrorNumber_OK)

		listRoomRet, err := uc.mCoreClient.ListRoom(context.TODO(), &core.ListRoomArgs{
			GameName: "", // show all room
		})

		So(err, assertions.ShouldBeNil)
		So(listRoomRet.Err, assertions.ShouldEqual, core.ErrorNumber_OK)
		So(len(listRoomRet.GetRoom()), assertions.ShouldEqual, 1)

		room := listRoomRet.GetRoom()[0]

		So(len(room.UserName), assertions.ShouldEqual, 0)

		joinRoomRet, err := uc.mCoreClient.JoinRoom(context.TODO(), &core.JoinRoomArgs{
			Token:  token,
			RoomId: room.GetRoomId(),
		})

		So(err, assertions.ShouldBeNil)
		So(joinRoomRet.Err, assertions.ShouldEqual, core.ErrorNumber_OK)

		listRoomRet2, err := uc.mCoreClient.ListRoom(context.TODO(), &core.ListRoomArgs{
			GameName: "", // show all room
		})

		So(err, assertions.ShouldBeNil)
		So(listRoomRet2.Err, assertions.ShouldEqual, core.ErrorNumber_OK)
		So(len(listRoomRet.GetRoom()), assertions.ShouldEqual, 1)

		room2 := listRoomRet2.GetRoom()[0]

		So(len(room2.UserName), assertions.ShouldEqual, 1)
		So(room2.UserName[0], assertions.ShouldEqual, userName)
	})
}
