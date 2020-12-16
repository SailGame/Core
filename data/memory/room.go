package memory

import (
	d "github.com/SailGame/Core/data"
)

type Room struct {
	mGameName string
	mUsers []d.User
	mProvider interface{}
	mState d.RoomState
}

func NewRoom() (Room){
	r := Room{}
	r.mUsers = make([]d.User, 0)
	r.mProvider = nil
	r.mState = d.Preparing
	return r
}

func (r Room) GetGameName() (string){
	return r.mGameName
}
func (r Room) SetGameName(name string){
	r.mGameName = name
}
func (r Room) GetUsers() ([]d.User){
	return r.mUsers
}
func (r Room) SetProvider(provider interface{}){
	r.mProvider = provider
}
func (r Room) GetProvider() (interface{}){
	return r.mProvider
}
func (r Room) GetState() (d.RoomState){
	return r.mState
}

func (r Room) UserJoin(user d.User) (error){
	r.mUsers = append(r.mUsers, user)
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