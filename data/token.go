package data

type Token interface {
	GetKey() string
	GetUser() User
	GetUserName() string
}
