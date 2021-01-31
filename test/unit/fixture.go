package unit

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/SailGame/Core/data/memory"
	cpb "github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/pb/core/mocks"
	"github.com/SailGame/Core/server"
	"github.com/golang/mock/gomock"
)

type fixture struct {
	ctrl        *gomock.Controller
	coreServer  *server.CoreServer
	userNum     int
	providerNum int
}

type provider struct {
	mId                 string
	mMockProviderServer *mocks.MockGameCore_ProviderServer
	mSendMsgCh          chan *cpb.ProviderMsg
}

type user struct {
	mMockUserServer *mocks.MockGameCore_ListenServer
	mUserName       string
	mToken          string
	mRoomId         int32
	mState          bool
}

func newFixture(t *testing.T) *fixture {
	coreServer, _ := server.NewCoreServer(&server.CoreServerConfig{
		MStorage: memory.NewStorage(),
	})
	f := &fixture{
		ctrl:       gomock.NewController(t),
		coreServer: coreServer,
		userNum:    0,
	}
	return f
}

func (f *fixture) done() {
	f.ctrl.Finish()
}

func (f *fixture) newMockProvider() *provider {
	p := &provider{
		mId:                 strconv.Itoa(f.providerNum),
		mMockProviderServer: mocks.NewMockGameCore_ProviderServer(f.ctrl),
		mSendMsgCh:          make(chan *cpb.ProviderMsg),
	}
	f.providerNum = f.providerNum + 1

	p.mMockProviderServer.EXPECT().Recv().AnyTimes().DoAndReturn(func() (*cpb.ProviderMsg, error) {
		msg := <-p.mSendMsgCh
		return msg, nil
	})
	go f.coreServer.Provider(p.mMockProviderServer)
	return p
}

func (f *fixture) newMockUser() *user {
	u := &user{
		mMockUserServer: mocks.NewMockGameCore_ListenServer(f.ctrl),
		mUserName:       strconv.Itoa(f.userNum),
		mToken:          "",
		mRoomId:         0,
		mState:          false,
	}
	f.userNum = f.userNum + 1

	loginRet, err := f.coreServer.Login(context.TODO(), &cpb.LoginArgs{
		UserName: u.mUserName,
	})
	if err != nil {
		log.Fatal(err)
	}

	u.mToken = loginRet.GetToken()

	go f.coreServer.Listen(&cpb.ListenArgs{
		Token: u.mToken,
	}, u.mMockUserServer)

	return u
}

func (f *fixture) newRoom(u *user) int32 {
	ret, err := f.coreServer.CreateRoom(context.TODO(), &cpb.CreateRoomArgs{
		Token: u.mToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	return ret.GetRoomId()
}

func (f *fixture) joinRoom(roomId int32, u *user) cpb.ErrorNumber {
	ret, err := f.coreServer.JoinRoom(context.TODO(), &cpb.JoinRoomArgs{
		Token:  u.mToken,
		RoomId: roomId,
	})
	if err != nil {
		log.Fatal(err)
	}
	u.mRoomId = roomId
	return ret.GetErr()
}

func (f *fixture) controlRoom(u *user, gameName string) cpb.ErrorNumber {
	ret, err := f.coreServer.ControlRoom(context.TODO(), &cpb.ControlRoomArgs{
		Token:    u.mToken,
		RoomId:   u.mRoomId,
		GameName: gameName,
	})
	if err != nil {
		log.Fatal(err)
	}
	return ret.GetErr()
}

func (f *fixture) userReady(u *user) cpb.ErrorNumber {
	ret, err := f.coreServer.OperationInRoom(context.TODO(), &cpb.OperationInRoomArgs{
		Token: u.mToken,
		RoomOperation: &cpb.OperationInRoomArgs_Ready{
			Ready: cpb.Ready_READY,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return ret.GetErr()
}

func (f *fixture) newUsersAndRoom(userNum int) (users []*user, roomID int32, p *provider) {
	if userNum < 0 {
		log.Fatalf("newUsersAndRoom: userNum < 0")
	}
	testGame := "testGame"
	p = f.newMockProvider()

	users = make([]*user, 0)
	for i := 0; i < userNum; i++ {
		u := f.newMockUser()
		u.mMockUserServer.EXPECT().Send(gomock.Any()).AnyTimes()
		users = append(users, u)
	}
	roomID = f.newRoom(users[0])

	// register provider
	p.mMockProviderServer.EXPECT().Send(gomock.Any()).Times(1)
	p.send(&cpb.ProviderMsg{
		Msg: &cpb.ProviderMsg_RegisterArgs{
			RegisterArgs: &cpb.RegisterArgs{
				Id:       p.mId,
				GameName: testGame,
				GameSetting: &cpb.GameSetting{
					MaxUsers: int32(userNum),
					MinUsers: int32(userNum),
				},
			},
		},
	})

	for _, u := range users {
		if f.joinRoom(roomID, u) != cpb.ErrorNumber_OK {
			log.Fatal("newUsersAndRoom: joinRoom fail")
		}
	}
	if f.controlRoom(users[0], testGame) != cpb.ErrorNumber_OK {
		log.Fatalf("newUsersAndRoom: controlRoom fail")
	}
	return
}


func (p *provider) send(msg *cpb.ProviderMsg) {
	p.mSendMsgCh <- msg
	time.Sleep(500 * time.Microsecond)
}