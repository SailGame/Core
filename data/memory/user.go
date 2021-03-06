package memory

import (
	"errors"
	"sync"

	d "github.com/SailGame/Core/data"
)

type User struct {
	mMutex       sync.Locker
	mUserName    string
	mPasswd      string
	mDisplayName string

	mConn        interface{}
	mRoom        d.Room
	mTemporaryID uint32
}

func NewUser(userName string, passwd string) *User {
	u := &User{
		mUserName: userName,
		mPasswd:   passwd,
		mMutex:    &sync.Mutex{},
	}
	return u
}

func (u *User) Lock() {
	u.mMutex.Lock()
}

func (u *User) Unlock() {
	u.mMutex.Unlock()
}

func (u *User) GetUserName() string {
	return u.mUserName
}

func (u *User) GetDisplayName() string {
	if u.mDisplayName == "" {
		return u.mUserName
	} else {
		return u.mDisplayName
	}
}

func (u *User) SetDisplayName(displayName string) error {
	u.mDisplayName = displayName
	return nil
}

func (u *User) SetPasswd(oldPasswd string, newPasswd string) error {
	if u.mPasswd == "" || u.mPasswd == oldPasswd {
		u.mPasswd = newPasswd
	} else {
		return errors.New("old passwd mismatch")
	}
	return nil
}

func (u *User) GetConn() (interface{}, error) {
	if u.mConn == nil {
		return nil, errors.New("No live connection")
	}
	return u.mConn, nil
}

func (u *User) SetConn(conn interface{}) {
	u.mConn = conn
}

func (u *User) GetRoom() (d.Room, error) {
	if u.mRoom != nil {
		return u.mRoom, nil
	} else {
		return nil, errors.New("Not in a room")
	}
}

func (u *User) SetRoom(room d.Room) error {
	u.mRoom = room
	return nil
}

func (u *User) GetTemporaryID() uint32 {
	return u.mTemporaryID
}

func (u *User) SetTemporaryID(tid uint32) {
	u.mTemporaryID = tid
}
