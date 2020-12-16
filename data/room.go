package data

type RoomState int32

const (
    Preparing       RoomState = 0
    Playing         RoomState = 1
)

type Room interface {
	GetGameName() (string)
	SetGameName(string)
	GetUsers() ([]User)
	SetProvider(interface{})
	GetProvider() (interface{})
	GetState() (RoomState)

	UserJoin(User) (error)
	UserReady(User) (error)
	UserExit(User) (error)
}