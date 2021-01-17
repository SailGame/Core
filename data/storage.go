package data

// Storage is an abstract layer between data user and data store
// all the type is interface
type Storage interface {
	CreateRoom() (Room, error)
	GetRooms() ([]Room)
	GetRoomsWithFilter(func(Room) bool) ([]Room)
	FindRoom(roomID int32) (Room, error)
	DelRoom(roomID int32) (error)

	IsUserExist(userName string) bool
	CreateUser(userName string, passwd string) (error)
	GetUsers() ([]User)
	FindUser(userName string, passwd string) (User, error)
	DelUser(userName string) (error)

	CreateToken(user User) (Token, error)
	FindToken(key string) (Token, error)
	DelToken(key string) (error)

	RegisterProvider(Provider) (error)
	GetProviders() ([]Provider)
	FindProvider(providerID string) (Provider, error)
	FindProviderByGame(gameName string) ([]Provider)
	UnRegisterProvider(Provider) (error)
}