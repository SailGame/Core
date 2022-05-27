package system

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/SailGame/Core/data/memory"
	cpb "github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/server"
)

func TestGameStart(t *testing.T) {
	Convey("Game Start", t, func() {
		f := newFixture()
		f.init(&server.CoreServerConfig{
			MStorage: memory.NewStorage(),
		})
		gameName := "testGame"

		uc, _, token, roomId := buildOneUserAndOneRoom(f)
		uc.listenToCore(token)

		gameSetting, err := ptypes.MarshalAny(&cpb.Account{
			UserName: "GameSetting",
			Points:   99,
		})

		So(err, ShouldBeNil)

		controlRoomRet, err := uc.mCoreClient.ControlRoom(context.TODO(), &cpb.ControlRoomArgs{
			Token:    token,
			RoomId:   roomId,
			GameName: gameName,
			Custom:   gameSetting,
		})

		So(err, ShouldBeNil)
		So(controlRoomRet.Err, ShouldEqual, cpb.ErrorNumber_ControlRoom_RequiredProviderNotExist)

		pcId := "testProvider"
		pc := f.newProviderClient()
		So(pc.connectToCore(), ShouldBeNil)

		err = pc.mProviderClient.Send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_RegisterArgs{
				&cpb.RegisterArgs{
					Id:       pcId,
					GameName: gameName,
				},
			},
		})

		So(err, ShouldBeNil)
		time.Sleep(500 * time.Millisecond)

		controlRoomRet, err = uc.mCoreClient.ControlRoom(context.TODO(), &cpb.ControlRoomArgs{
			Token:    token,
			RoomId:   roomId,
			GameName: gameName,
			Custom:   gameSetting,
		})

		So(err, ShouldBeNil)
		So(controlRoomRet.Err, ShouldEqual, cpb.ErrorNumber_OK)

		qryRoomRet, err := uc.mCoreClient.QueryRoom(context.TODO(), &cpb.QueryRoomArgs{
			Token:  token,
			RoomId: roomId,
		})

		So(err, ShouldBeNil)
		So(qryRoomRet.Err, ShouldEqual, cpb.ErrorNumber_OK)
		So(qryRoomRet.Room.RoomId, ShouldEqual, roomId)
		So(qryRoomRet.Room.GameName, ShouldEqual, gameName)

		unMarshalGameSetting := &cpb.Account{}
		err = ptypes.UnmarshalAny(qryRoomRet.GetRoom().GetGameSetting(), unMarshalGameSetting)
		So(err, ShouldBeNil)
		So(unMarshalGameSetting.UserName, ShouldEqual, "GameSetting")
		So(unMarshalGameSetting.Points, ShouldEqual, 99)

		uc.mCoreClient.OperationInRoom(context.TODO(), &cpb.OperationInRoomArgs{
			Token: token,
			RoomOperation: &cpb.OperationInRoomArgs_Ready{
				Ready: cpb.Ready_READY,
			},
		})

		// all user is ready, user and provider should receive the start signal
		var wg sync.WaitGroup
		wg.Add(2)
		go Convey("uc recv", t, func() {
			JoinMsg, err := uc.mLisClient.Recv()
			So(err, ShouldBeNil)
			So(JoinMsg.GetRoomDetails(), ShouldNotBeNil)
			ReadyMsg, err := uc.mLisClient.Recv()
			So(err, ShouldBeNil)
			So(ReadyMsg.GetRoomDetails(), ShouldNotBeNil)
			msg, err := uc.mLisClient.Recv()

			t.Logf("TestUser receive start game")
			So(err, ShouldBeNil)
			So(msg.GetCustom(), ShouldNotBeNil)
			wg.Done()
		})

		go Convey("pc recv", t, func() {
			regRetMsg, err := pc.mProviderClient.Recv()
			So(err, ShouldBeNil)
			So(regRetMsg.GetRegisterRet(), ShouldNotBeNil)

			msg, err := pc.mProviderClient.Recv()
			So(err, ShouldBeNil)
			So(msg.GetStartGameArgs(), ShouldNotBeNil)
			err = ptypes.UnmarshalAny(msg.GetStartGameArgs().GetCustom(), unMarshalGameSetting)
			So(err, ShouldBeNil)
			So(unMarshalGameSetting.UserName, ShouldEqual, "GameSetting")
			So(unMarshalGameSetting.Points, ShouldEqual, 99)

			t.Logf("TestProvider send start to TestUser")
			err = pc.mProviderClient.Send(&cpb.ProviderMsg{
				Msg: &cpb.ProviderMsg_NotifyMsgArgs{
					&cpb.NotifyMsgArgs{
						RoomId: msg.GetStartGameArgs().GetRoomId(),
						UserId: 0, // broadcast
						Custom: nil,
					},
				},
			})
			So(err, ShouldBeNil)
			wg.Done()
		})

		ch := make(chan int)
		go func() {
			wg.Wait()
			ch <- 1
		}()

		select {
		case <-ch:
			return
		case <-time.After(5 * time.Second):
			t.Fatalf("timeout")
		}
	})
}
