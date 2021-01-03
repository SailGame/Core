package data

type RoomState int32

const (
    RoomState_PREPARING       RoomState = 0
    RoomState_PLAYING         RoomState = 1
)

type UserState int32

const (
	UserState_ERROR			  UserState = 0
	UserState_PREPARING       UserState = 1
	UserState_READY           UserState = 2
	UserState_PLAYING	      UserState = 3
	// will be deleted after game over
	UserState_EXITED          UserState = 4
)

type Room interface {
	// concurrency control
	Lock()
	Unlock()
	GetRoomID() (int32)
	GetGameName() (string)
	GetUsers() ([]User)
	SetProvider(Provider)
	GetProvider() (Provider)
	GetState() (RoomState)
	GetUserState(User) (UserState, error)

	UserJoin(User) (error)
	UserReady(User, bool) (error)
	UserExit(User) (error)

	Restart() (error)
}