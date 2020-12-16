package data

type RoomState int32

const (
    Preparing       RoomState = 0
    Playing         RoomState = 1
)

type Room interface {
	GetGameName() (string, error)
	SetGameName(string)
	GetUsers() ([]User, error)
	SetProvider(interface{})
	GetProvider() (interface{}, error)
	GetState() (RoomState)

	UserJoin(User) (error)
	UserReady(User) (error)
	UserExit(User) (error)
}