package memory

import (
	d "github.com/SailGame/Core/data"
)

type Room struct {
	mRoomID int32
	mGameName string
	mUsers map[string]d.User
	mProvider d.Provider
	mState d.RoomState
}

func NewRoom(ID int32) (Room){
	r := Room{}
	r.mRoomID = ID
	r.mUsers = make(map[string]d.User)
	r.mProvider = nil
	r.mState = d.Preparing
	return r
}

func (r Room) GetRoomID() (int32){
	return r.mRoomID
}

func (r Room) GetGameName() (string){
	return r.mGameName
}

func (r Room) SetGameName(name string){
	r.mGameName = name
}

func (r Room) GetUsers() ([]d.User){
	ret := make([]d.User, len(r.mUsers))
	for _, v := range r.mUsers {
		ret = append(ret, v)
	}
	return ret
}

func (r Room) SetProvider(provider d.Provider){
	r.mProvider = provider
}

func (r Room) GetProvider() (d.Provider){
	return r.mProvider
}

func (r Room) GetState() (d.RoomState){
	return r.mState
}

func (r Room) UserJoin(user d.User) (error){
	// TODO: check room capacity
	r.mUsers[user.GetUserName()] = user
	user.SetRoom(r)
	return nil
}

func (r Room) UserReady(d.User) (error){
	// TODO: if all users are ready, start the game
	return nil
}
func (r Room) UserExit(d.User) (error){
	// TODO: delete the user
	return nil
}