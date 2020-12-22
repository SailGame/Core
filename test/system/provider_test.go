package system

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
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

		controlRoomRet, err := uc.mCoreClient.ControlRoom(context.TODO(), &cpb.ControlRoomArgs{
			Token: token,
			RoomId: roomId,
			GameName: gameName,
		})
		So(err, assertions.ShouldBeNil)
		So(controlRoomRet.Errno, assertions.ShouldEqual, cpb.ErrorNumber_ControlRoom_RequiredProviderNotExist)

		pcId := "testProvider"
		pc := f.newProviderClient()
		err = pc.connectToCore()
		So(err, assertions.ShouldBeNil)

		err = pc.mProviderClient.Send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_RegisterArgs{
				&cpb.RegisterArgs{
					Id: pcId,
					GameName: gameName,
				},
			},
		})

		So(err, assertions.ShouldBeNil)
		time.Sleep(500 * time.Millisecond)

		controlRoomRet, err = uc.mCoreClient.ControlRoom(context.TODO(), &cpb.ControlRoomArgs{
			Token: token,
			RoomId: roomId,
			GameName: gameName,
		})

		So(err, assertions.ShouldBeNil)
		So(controlRoomRet.Errno, assertions.ShouldEqual, cpb.ErrorNumber_OK)

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
			msg, err := uc.mLisClient.Recv()

			So(err, assertions.ShouldBeNil)
			So(msg.GetCustom(), assertions.ShouldNotBeNil)
			wg.Done()
		})

		go Convey("pc recv", t, func() {
			msg, err := pc.mProviderClient.Recv()

			So(err, assertions.ShouldBeNil)
			So(msg.GetStartGameArgs(), assertions.ShouldNotBeNil)
			wg.Done()
		})

		var ch chan interface{}
		go func() {
			wg.Wait()
			ch <- nil
		}()

		select {
		case <-ch:
			return
		case <-time.After(1 * time.Second):
			t.Fatalf("timeout")
		}
	})
}