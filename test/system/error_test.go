package system

import (
	"context"
	"testing"

	"github.com/smartystreets/assertions"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/SailGame/Core/data/memory"
	cpb "github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/server"
)

func TestClientDisconnect(t *testing.T) {
	Convey("Client Disconnect", t, func() {
		f := newFixture()
		f.init(&server.CoreServerConfig{
			MStorage: memory.NewStorage(),
		})
		gameName := "testGame"

		uc, _, token, roomId := buildOneUserAndOneRoom(f)
		uc.listenToCore(token)

		controlRoomRet, err := uc.mCoreClient.JoinRoom(context.TODO(), &cpb.JoinRoomArgs{
			Token:    token,
			RoomId:   roomId,
		})
		So(err, assertions.ShouldBeNil)
		So(controlRoomRet.Err, assertions.ShouldEqual, cpb.ErrorNumber_OK)

		pcId := "testProvider"
		pc := f.newProviderClient()
		So(pc.connectToCore(), assertions.ShouldBeNil)

		err = pc.mProviderClient.Send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_RegisterArgs{
				&cpb.RegisterArgs{
					Id:       pcId,
					GameName: gameName,
				},
			},
		})

		So(err, assertions.ShouldBeNil)

		uc.mClose()
		_, err = uc.mLisClient.Recv()
		So(err, assertions.ShouldNotBeNil)

		// all user is ready, user and provider should receive the start signal
		// but user disconnects from core, core should discard the msg to that user
		t.Logf("TestProvider send msg to TestUser")
		err = pc.mProviderClient.Send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_NotifyMsgArgs{
				&cpb.NotifyMsgArgs{
					RoomId: roomId,
					UserId: 0, // broadcast
					Custom: nil,
				},
			},
		})
		So(err, assertions.ShouldBeNil)

		_, err = uc.mLisClient.Recv()
		So(err, assertions.ShouldNotBeNil)
	})
}
