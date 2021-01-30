package unit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	cpb "github.com/SailGame/Core/pb/core"
	"github.com/golang/mock/gomock"
)

func TestMaxUser(t *testing.T) {
	Convey("Test Max User", t, func() {
		testGame := "testGame"
		f := newFixture(t)
		p := f.newMockProvider()
		u1 := f.newMockUser()
		u2 := f.newMockUser()

		room := f.newRoom(u1)

		// register provider
		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		p.mSendMsgCh <- &cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_RegisterArgs{
				RegisterArgs: &cpb.RegisterArgs{
					Id: p.mId,
					GameName: testGame,
					GameSetting: &cpb.GameSetting{
						MaxUsers: 1,
						MinUsers: 1,
					},
				},
			},
		}

		So(f.joinRoom(room, u1), ShouldEqual, cpb.ErrorNumber_OK)
		So(f.joinRoom(room, u2), ShouldEqual, cpb.ErrorNumber_JoinRoom_FullRoom)
	})
}

func TestMinUser(t *testing.T) {
	Convey("Test Min User", t, func() {
		testGame := "testGame"
		f := newFixture(t)
		p := f.newMockProvider()
		u1 := f.newMockUser()
		u2 := f.newMockUser()

		room := f.newRoom(u1)

		// register provider
		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		p.mSendMsgCh <- &cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_RegisterArgs{
				RegisterArgs: &cpb.RegisterArgs{
					Id: p.mId,
					GameName: testGame,
					GameSetting: &cpb.GameSetting{
						MaxUsers: 5,
						MinUsers: 2,
					},
				},
			},
		}

		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(0)
		So(f.joinRoom(room, u1), ShouldEqual, cpb.ErrorNumber_OK)
		So(f.userReady(u1), ShouldEqual, cpb.ErrorNumber_OK)

		So(f.joinRoom(room, u2), ShouldEqual, cpb.ErrorNumber_OK)
		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		So(f.userReady(u2), ShouldEqual, cpb.ErrorNumber_OK)
	})
}
