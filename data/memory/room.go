package memory

import (
	"errors"
	"sync"

	d "github.com/SailGame/Core/data"
)

type userState int32

const (
	preparing       userState = 0
	ready           userState = 1
)
type Room struct {
	mRoomID int32
	mGameName string
	mUsers map[string]d.User
	mUserState map[string]userState
	mProvider d.Provider
	mState d.RoomState
	mMutex sync.Locker
}

func NewRoom(ID int32) (*Room){
	r := Room{}
	r.mRoomID = ID
	r.mUsers = make(map[string]d.User)
	r.mUserState = make(map[string]userState)
	r.mProvider = nil
	r.mState = d.Preparing
	r.mMutex = &sync.Mutex{}
	return &r
}

func (r *Room) GetRoomID() (int32){
	return r.mRoomID
}

func (r *Room) GetGameName() (string){
	if(r.mProvider == nil){
		return ""
	}
	return r.mProvider.GetGameName()
}

func (r *Room) GetUsers() ([]d.User){
	ret := make([]d.User, 0, len(r.mUsers))
	r.mMutex.Lock()
	defer r.mMutex.Unlock()
	for _, v := range r.mUsers {
		ret = append(ret, v)
	}
	return ret
}

func (r *Room) SetProvider(provider d.Provider){
	if(r.mProvider == provider){
		return
	}
	if(r.mProvider != nil){
		r.mProvider.DelRoom(r)
	}
	r.mProvider = provider
	r.mGameName = provider.GetGameName()
	provider.AddRoom(r)
}

func (r *Room) GetProvider() (d.Provider){
	return r.mProvider
}

func (r *Room) GetState() (d.RoomState){
	return r.mState
}

func (r *Room) UserJoin(user d.User) (error){
	// TODO: check room capacity
	r.mMutex.Lock()
	defer r.mMutex.Unlock()
	r.mUsers[user.GetUserName()] = user
	r.mUserState[user.GetUserName()] = preparing
	user.SetRoom(r)
	return nil
}

func (r *Room) UserReady(user d.User, ok bool) (error){
	r.mMutex.Lock()
	defer r.mMutex.Unlock()
	if(r.mState == d.Playing){
		return errors.New("Not support change user state when game is playing")
	}
	if(ok){
		r.mUserState[user.GetUserName()] = ready
	}else {
		r.mUserState[user.GetUserName()] = preparing
	}
	for _, v := range r.mUserState {
		if(v != ready){
			return nil
		}
	}
	r.mState = d.Playing;
	return nil
}

func (r *Room) UserExit(user d.User) (error){
	r.mMutex.Lock()
	defer r.mMutex.Unlock()
	delete(r.mUsers, user.GetUserName())
	delete(r.mUserState, user.GetUserName())
	return nil
}

func (r *Room) Restart() (error){
	r.mMutex.Lock()
	defer r.mMutex.Unlock()
	for un, _ := range r.mUserState {
		r.mUserState[un] = preparing
	}
	return nil
}