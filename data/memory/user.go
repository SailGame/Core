package memory

import (
	"errors"

	d "github.com/SailGame/Core/data"
)

type User struct {
	mUserName string
	mPasswd string
	mDisplayName string
	mRoom d.Room
}

func NewUser(userName string, passwd string) (User){
	u := User{}
	u.mUserName = userName
	u.mPasswd = passwd
	return u
}

func (u User) GetUserName() (string){
	return u.mUserName
}

func (u User) GetDisplayName() (string){
	if(u.mDisplayName == ""){
		return u.mUserName
	}else{
		return u.mDisplayName
	}
}

func (u User) SetDisplayName(displayName string) (error){
	u.mDisplayName = displayName
	return nil
}

func (u User) SetPasswd(oldPasswd string, newPasswd string) (error){
	if(u.mPasswd == "" || u.mPasswd == oldPasswd){
		u.mPasswd = newPasswd
	}else{
		return errors.New("old passwd mismatch")
	}
	return nil
}

func (u User) GetRoom() (d.Room, error){
	if(u.mRoom != nil){
		return u.mRoom, nil
	}else{
		return nil, errors.New("Not in a room")
	}
}
func (u User) SetRoom(room d.Room) (error){
	if(u.mRoom != nil && u.mRoom != room){
		u.mRoom.UserExit(u)
	}
	u.mRoom = room
	return nil
}