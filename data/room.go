package data

type RoomState int32

const (
    Preparing       RoomState = 0
    Playing         RoomState = 1
)

type Room interface {
	GetRoomID() (int32)
	GetGameName() (string)
	SetGameName(string)
	GetUsers() ([]User)
	SetProvider(Provider)
	GetProvider() (Provider)
	GetState() (RoomState)

	UserJoin(User) (error)
	UserReady(User) (error)
	UserExit(User) (error)
}