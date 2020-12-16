package memory

import (
	"errors"
	"log"

	d "github.com/SailGame/Core/data"
)

type Storage struct {
	mRooms map[int]Room
	mUsers map[string]User
	mTokens map[string]Token
	mProviders map[string]interface{}
}

func NewStorage() (Storage){
	storage := Storage{}
	storage.mRooms = make(map[int]Room)
	storage.mUsers = make(map[string]User)
	storage.mTokens = make(map[string]Token)
	storage.mProviders = make(map[string]interface{})
	return storage
}

func (s Storage) CreateRoom() (error){
	// TODO: if the room is released, the length is not a stable id
	s.mRooms[len(s.mRooms)] = NewRoom()
	return nil
}

func (s Storage) FindRoom(roomID int) (d.Room, error){
	room, ok := s.mRooms[roomID]
	if(ok){
		return room, nil
	} else{
		log.Printf("Room (%d) not exist", roomID)
		return nil, errors.New("Room not exist");
	}
}

// func (s *Storage) DelRoom(roomID int) (error){

// }

// func (s *Storage) CreateUser(userName string, passwd string) (error){

// }
// func (s *Storage) FindUser(userName string, passwd string) (d.User, error){

// }
// func (s *Storage) DelUser(userName string) (error){

// }

// func (s *Storage) CreateToken(userID int) (error){

// }
// func (s *Storage) FindToken(token string) (d.Token, error){

// }
// func (s *Storage) DelToken(token string) (error){

// }
// func (s *Storage) RegisterProvider(providerID string, provider interface{}) (error){

// }
// func (s *Storage) FindProvider(providerID string) (interface{}, error){

// }
// func (s *Storage) UnRegisterProvider(providerID string) (error){

// }

type User struct {

}

type Token struct {

}