package data

type Token interface {
	GetToken() (string)
	GetUserName() (string)
}