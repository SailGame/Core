package data

type User interface {
	GetUserName() (string)
	GetDisplayName() (string)
	SetDisplayName() (string, error)
	SetPasswd(oldPasswd string, newPasswd string) (error)
}