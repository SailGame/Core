package data

type Provider interface {
	GetID() (string)
	GetGameName() (string)
	GetRooms() ([]Room)

	AddRoom(Room) (error)
	DelRoom(Room) (error)
}