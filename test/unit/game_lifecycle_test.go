package unit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	cpb "github.com/SailGame/Core/pb/core"
	"github.com/golang/mock/gomock"
)

const (
	userNum int = 4
)

// ready, start, close
func TestCompleteLifecycle(t *testing.T) {
	Convey("Test Complete Lifecycle", t, func() {
		f := newFixture(t)
		defer f.done()
		users, roomID, p := f.newUsersAndRoom(userNum)

		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		for _, u := range users {
			So(f.userReady(u), ShouldEqual, cpb.ErrorNumber_OK)
		}

		p.send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_CloseGameArgs{
				CloseGameArgs: &cpb.CloseGameArgs{
					RoomId: roomID,
				},
			},})

		p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
		for _, u := range users {
			So(f.userReady(u), ShouldEqual, cpb.ErrorNumber_OK)
		}
	})
}
