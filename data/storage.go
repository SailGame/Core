package data

// Storage is an abstract layer between data user and data store
// all the type is interface
type Storage interface {
	CreateRoom() (Room, error)
	GetRooms() ([]Room)
	FindRoom(roomID int32) (Room, error)
	DelRoom(roomID int32) (error)

	CreateUser(userName string, passwd string) (error)
	GetUsers() ([]User)
	FindUser(userName string, passwd string) (User, error)
	DelUser(userName string) (error)

	CreateToken(user User) (Token, error)
	FindToken(key string) (Token, error)
	DelToken(key string) (error)

	RegisterProvider(providerID string, provider Provider) (error)
	GetProviders() ([]Provider)
	FindProvider(providerID string) (Provider, error)
	FindProviderByGame(gameName string) ([]Provider)
	UnRegisterProvider(providerID string) (error)
}