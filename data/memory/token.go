package memory

import (
	d "github.com/SailGame/Core/data"
)

type Token struct {
	mKey  string
	mUser d.User
}

func (t Token) GetKey() string {
	return t.mKey
}

func (t Token) GetUser() d.User {
	return t.mUser
}
func (t Token) GetUserName() string {
	return t.mUser.GetUserName()
}
