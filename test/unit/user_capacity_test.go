package unit

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	cpb "github.com/SailGame/Core/pb/core"
	"github.com/golang/mock/gomock"
)

const (
	testGame string = "testGame"
)

func TestMaxUser(t *testing.T) {
	Convey("Test Max User", t, func() {
		f := newFixture(t)
		defer f.done()
		p := f.newMockProvider()
		u1 := f.newMockUser()
		u2 := f.newMockUser()

		room := f.newRoom(u1)

		// register provider
		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		p.send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_RegisterArgs{
				RegisterArgs: &cpb.RegisterArgs{
					Id:       p.mId,
					GameName: testGame,
					GameSetting: &cpb.GameSetting{
						MaxUsers: 1,
						MinUsers: 1,
					},
				},
			},
		})
		time.Sleep(500 * time.Microsecond)
		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.joinRoom(room, u1), ShouldEqual, cpb.ErrorNumber_OK)
		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.controlRoom(u1, testGame), ShouldEqual, cpb.ErrorNumber_OK)

		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(0)
		u2.mMockUserServer.EXPECT().Send(gomock.Any()).Times(0)
		So(f.joinRoom(room, u2), ShouldEqual, cpb.ErrorNumber_JoinRoom_FullRoom)
	})
}

func TestMinUser(t *testing.T) {
	Convey("Test Min User", t, func() {
		f := newFixture(t)
		defer f.done()
		p := f.newMockProvider()
		u1 := f.newMockUser()
		u2 := f.newMockUser()

		room := f.newRoom(u1)

		// register provider
		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		p.send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_RegisterArgs{
				RegisterArgs: &cpb.RegisterArgs{
					Id:       p.mId,
					GameName: testGame,
					GameSetting: &cpb.GameSetting{
						MaxUsers: 5,
						MinUsers: 2,
					},
				},
			},
		})
		time.Sleep(500 * time.Microsecond)

		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(0)
		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.joinRoom(room, u1), ShouldEqual, cpb.ErrorNumber_OK)
		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.controlRoom(u1, testGame), ShouldEqual, cpb.ErrorNumber_OK)
		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.userReady(u1), ShouldEqual, cpb.ErrorNumber_OK)

		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		u2.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.joinRoom(room, u2), ShouldEqual, cpb.ErrorNumber_OK)
		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		u1.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		u2.mMockUserServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.userReady(u2), ShouldEqual, cpb.ErrorNumber_OK)
	})
}
