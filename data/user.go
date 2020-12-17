package data

type User interface {
	// user info
	GetUserName() (string)
	GetDisplayName() (string)
	SetDisplayName(string) (error)
	SetPasswd(oldPasswd string, newPasswd string) (error)

	// game
	GetRoom() (Room, error)
	SetRoom(Room) (error)
}