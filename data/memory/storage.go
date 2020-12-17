package memory

import (
	"errors"
	"log"

	d "github.com/SailGame/Core/data"
	"github.com/go-basic/uuid"
)

type Storage struct {
	mRooms map[int32]Room
	mUsers map[string]User
	mTokens map[string]Token
	mProviders map[string]d.Provider
}

func NewStorage() (Storage){
	storage := Storage{}
	storage.mRooms = make(map[int32]Room)
	storage.mUsers = make(map[string]User)
	storage.mTokens = make(map[string]Token)
	storage.mProviders = make(map[string]d.Provider)
	return storage
}

func (s Storage) CreateRoom() (d.Room, error){
	// TODO: if the room is released, the length is not a stable id
	newRoom := NewRoom()
	s.mRooms[int32(len(s.mRooms))] = newRoom
	return newRoom, nil
}

func (s Storage) FindRoom(roomID int32) (d.Room, error){
	room, ok := s.mRooms[roomID]
	if(ok){
		return room, nil
	} else{
		log.Printf("Room (%d) not exist", roomID)
		return nil, errors.New("Room not exist");
	}
}

func (s Storage) DelRoom(roomID int32) (error){
	_, ok := s.mRooms[roomID]
	if(!ok){
		return errors.New("No such room")
	}
	delete(s.mRooms, roomID)
	return nil
}

func (s Storage) CreateUser(userName string, passwd string) (error){
	_, ok := s.mUsers[userName]
	if(ok){
		return errors.New("UserName is occupied:" + userName)
	}
	s.mUsers[userName] = NewUser(userName, passwd)
	return nil
}

func (s Storage) FindUser(userName string, passwd string) (d.User, error){
	user, ok := s.mUsers[userName]
	if(ok){
		return nil, errors.New("No such user:" + userName)
	}
	return user, nil
}

func (s Storage) DelUser(userName string) (error){
	// TODO: is user playing?
	delete(s.mUsers, userName)
	return nil
}

func (s Storage) CreateToken(user d.User) (error){
	// TODO: clear old token?
	uuid := uuid.New()
	s.mTokens[uuid] = Token{mKey: uuid, mUser: user}
	return nil
}

func (s Storage) FindToken(key string) (d.Token, error){
	token, ok := s.mTokens[key]
	if(ok){
		return token, nil
	}else{
		return nil, errors.New("No such token:" + key)
	}
}

func (s Storage) DelToken(key string) (error){
	delete(s.mTokens, key)
	return nil
}

func (s Storage) RegisterProvider(providerID string, provider d.Provider) (error){
	s.mProviders[providerID] = provider
	return nil
}

func (s Storage) FindProvider(providerID string) (d.Provider, error){
	p, ok := s.mProviders[providerID]
	if(!ok){
		return nil, errors.New("No such provider:" + providerID)
	}
	return p, nil
}

func (s Storage) FindProviderByGame(gameName string) ([]d.Provider){
	ret := make([]d.Provider, 0)
	for _, v := range s.mProviders {
		if(v.GetGameName() == gameName){
			ret = append(ret, v)
		}
	}
	return ret
}

func (s Storage) UnRegisterProvider(providerID string) (error){
	delete(s.mProviders, providerID)
	return nil
}