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
	GetTemporaryID() (uint32)	// avoid exposing the username to provider
	SetTemporaryID(uint32)
}