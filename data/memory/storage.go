package memory

import (
	"errors"
	"log"
	"sync"

	d "github.com/SailGame/Core/data"
	"github.com/go-basic/uuid"
)

type Storage struct {
	mRooms map[int32]Room
	mUsers map[string]User
	mTokens map[string]Token
	mProviders map[string]d.Provider
	mMutex sync.Locker
}

func NewStorage() (Storage){
	storage := Storage{}
	storage.mRooms = make(map[int32]Room)
	storage.mUsers = make(map[string]User)
	storage.mTokens = make(map[string]Token)
	storage.mProviders = make(map[string]d.Provider)
	storage.mMutex = &sync.Mutex{}
	return storage
}

func (s Storage) CreateRoom() (d.Room, error){
	// TODO: if the room is released, the length is not a stable id
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	roomID := int32(len(s.mRooms));
	newRoom := NewRoom(roomID)
	s.mRooms[roomID] = newRoom
	return newRoom, nil
}

func (s Storage) GetRooms() ([]d.Room){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.Room, len(s.mRooms))
	for _, v := range s.mRooms {
		ret = append(ret, v)
	}
	return ret
}

func (s Storage) FindRoom(roomID int32) (d.Room, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	room, ok := s.mRooms[roomID]
	if(ok){
		return room, nil
	} else{
		log.Printf("Room (%d) not exist", roomID)
		return nil, errors.New("Room not exist");
	}
}

func (s Storage) DelRoom(roomID int32) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	_, ok := s.mRooms[roomID]
	if(!ok){
		return errors.New("No such room")
	}
	delete(s.mRooms, roomID)
	return nil
}

func (s Storage) CreateUser(userName string, passwd string) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	_, ok := s.mUsers[userName]
	if(ok){
		return errors.New("UserName is occupied:" + userName)
	}
	s.mUsers[userName] = NewUser(userName, passwd)
	return nil
}

func (s Storage) GetUsers() ([]d.User){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.User, len(s.mUsers))
	for _, v := range s.mUsers {
		ret = append(ret, v)
	}
	return ret
}

func (s Storage) FindUser(userName string, passwd string) (d.User, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	user, ok := s.mUsers[userName]
	if(ok){
		return nil, errors.New("No such user:" + userName)
	}
	return user, nil
}

func (s Storage) DelUser(userName string) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	// TODO: is user playing?
	delete(s.mUsers, userName)
	return nil
}

func (s Storage) CreateToken(user d.User) (Token, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	// TODO: clear old token?
	uuid := uuid.New()
	newToken := Token{mKey: uuid, mUser: user}
	s.mTokens[uuid] = newToken
	return newToken, nil
}

func (s Storage) FindToken(key string) (d.Token, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	token, ok := s.mTokens[key]
	if(ok){
		return token, nil
	}else{
		return nil, errors.New("No such token:" + key)
	}
}

func (s Storage) DelToken(key string) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	delete(s.mTokens, key)
	return nil
}

func (s Storage) RegisterProvider(providerID string, provider d.Provider) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	s.mProviders[providerID] = provider
	return nil
}

func (s Storage) GetProviders() ([]d.Provider){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.Provider, len(s.mProviders))
	for _, v := range s.mProviders {
		ret = append(ret, v)
	}
	return ret
}

func (s Storage) FindProvider(providerID string) (d.Provider, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	p, ok := s.mProviders[providerID]
	if(!ok){
		return nil, errors.New("No such provider:" + providerID)
	}
	return p, nil
}

func (s Storage) FindProviderByGame(gameName string) ([]d.Provider){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.Provider, len(s.mProviders))
	for _, v := range s.mProviders {
		if(v.GetGameName() == gameName){
			ret = append(ret, v)
		}
	}
	return ret
}

func (s Storage) UnRegisterProvider(providerID string) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	delete(s.mProviders, providerID)
	return nil
}