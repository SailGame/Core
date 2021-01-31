package memory

import (
	"errors"
	"fmt"
	"sync"

	d "github.com/SailGame/Core/data"
	"google.golang.org/protobuf/types/known/anypb"
)

func isFullRoom(cgs *d.CommonGameSetting, num int32) bool {
	return (cgs.MaxUser != 0 && num >= cgs.MaxUser)
}

func isAcceptableUserNum(cgs *d.CommonGameSetting, num int32) bool {
	return (cgs.MaxUser == 0 || num <= cgs.MaxUser) && cgs.MinUser <= num
}

type Room struct {
	mRoomID            int32
	mUsers             map[string]d.User
	mUserStates        map[string]d.UserState
	mProvider          d.Provider
	mState             d.RoomState
	mCustomGameSetting *anypb.Any
	mMutex             sync.Locker
}

func NewRoom(ID int32) *Room {
	r := &Room{
		mRoomID:            ID,
		mUsers:             make(map[string]d.User),
		mUserStates:        make(map[string]d.UserState),
		mProvider:          nil,
		mState:             d.RoomState_PREPARING,
		mCustomGameSetting: nil,
		mMutex:             &sync.Mutex{},
	}
	return r
}

func (r *Room) Lock() {
	r.mMutex.Lock()
}

func (r *Room) Unlock() {
	r.mMutex.Unlock()
}

func (r *Room) GetRoomID() int32 {
	return r.mRoomID
}

func (r *Room) GetGameName() string {
	if r.mProvider == nil {
		return ""
	}
	return r.mProvider.GetGameSetting().GameName
}

func (r *Room) GetCustomGameSetting() *anypb.Any {
	return r.mCustomGameSetting
}

func (r *Room) SetCustomGameSetting(setting *anypb.Any) {
	r.mCustomGameSetting = setting
}

func (r *Room) GetUsers() []d.User {
	ret := make([]d.User, 0, len(r.mUsers))
	for _, v := range r.mUsers {
		ret = append(ret, v)
	}
	return ret
}

func (r *Room) SetProvider(provider d.Provider) {
	if r.mProvider == provider {
		return
	}
	if r.mProvider != nil {
		r.mProvider.DelRoom(r)
	}
	r.mProvider = provider
	provider.AddRoom(r)
}

func (r *Room) GetProvider() d.Provider {
	return r.mProvider
}

func (r *Room) GetState() d.RoomState {
	return r.mState
}

func (r *Room) GetUserState(user d.User) (d.UserState, error) {
	state, ok := r.mUserStates[user.GetUserName()]
	if !ok {
		return d.UserState_ERROR, errors.New("No such user")
	}
	return state, nil
}

func (r *Room) UserJoin(user d.User) error {
	if r.mProvider != nil && isFullRoom(r.mProvider.GetGameSetting(), int32(len(r.GetUsers()))) {
		return errors.New("Full Room")
	}
	if r.mState == d.RoomState_PLAYING {
		_, ok := r.mUsers[user.GetUserName()]
		if ok && r.mUserStates[user.GetUserName()] == d.UserState_EXITED {
			r.mUserStates[user.GetUserName()] = d.UserState_PLAYING
		} else {
			return errors.New("Not support change user state when game is playing")
		}
	}
	r.mUsers[user.GetUserName()] = user
	r.mUserStates[user.GetUserName()] = d.UserState_PREPARING
	return nil
}

func (r *Room) UserReady(user d.User, isReady bool) error {
	_, ok := r.mUsers[user.GetUserName()]
	if !ok {
		return fmt.Errorf("No such user(%s) in room(%d)", user.GetUserName(), r.mRoomID)
	}
	if r.mState == d.RoomState_PLAYING {
		return errors.New("Not support change user state when game is playing")
	}
	if isReady {
		r.mUserStates[user.GetUserName()] = d.UserState_READY
	} else {
		r.mUserStates[user.GetUserName()] = d.UserState_PREPARING
	}

	if r.mProvider != nil && isAcceptableUserNum(r.mProvider.GetGameSetting(), int32(len(r.mUsers))) {
		for _, v := range r.mUserStates {
			if v != d.UserState_READY {
				return nil
			}
		}
		for k := range r.mUserStates {
			r.mUserStates[k] = d.UserState_PLAYING
		}
		r.mState = d.RoomState_PLAYING
	}

	return nil
}

func (r *Room) UserExit(user d.User) (bool, error) {
	_, ok := r.mUsers[user.GetUserName()]
	if !ok {
		return false, fmt.Errorf("No such user(%s) in room(%d)", user.GetUserName(), r.mRoomID)
	}
	if r.mState == d.RoomState_PLAYING {
		r.mUserStates[user.GetUserName()] = d.UserState_EXITED
		return false, nil
	} else {
		delete(r.mUsers, user.GetUserName())
		delete(r.mUserStates, user.GetUserName())
	}
	return true, nil
}

func (r *Room) Restart() error {
	for un, state := range r.mUserStates {
		if state == d.UserState_PLAYING {
			r.mUserStates[un] = d.UserState_PREPARING
		} else if state == d.UserState_EXITED {
			r.mUsers[un].SetRoom(nil)
			delete(r.mUsers, un)
			delete(r.mUserStates, un)
		}
	}
	r.mState = d.RoomState_PREPARING
	return nil
}
