package data

// Storage is an abstract layer between data user and data store
// all the type is interface
type Storage interface {
	CreateRoom() (Room, error)
	FindRoom(roomID int32) (Room, error)
	DelRoom(roomID int32) (error)

	CreateUser(userName string, passwd string) (error)
	FindUser(userName string, passwd string) (User, error)
	DelUser(userName string) (error)

	CreateToken(user User) (error)
	FindToken(key string) (Token, error)
	DelToken(key string) (error)

	RegisterProvider(providerID string, provider interface{}) (error)
	FindProvider(providerID string) (interface{}, error)
	FindProviderByGame(gameName string) ([]interface{}, error)
	UnRegisterProvider(providerID string) (error)
}