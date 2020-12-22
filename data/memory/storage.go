package memory

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"

	d "github.com/SailGame/Core/data"
	"github.com/go-basic/uuid"
)

type Storage struct {
	mRooms map[int32]*Room
	mUsers map[string]*User
	mTokens map[string]*Token
	mProviders map[string]d.Provider
	mMutex sync.Locker
	mIdentityID int32
}

func NewStorage() (*Storage){
	storage := Storage{}
	storage.mRooms = make(map[int32]*Room)
	storage.mUsers = make(map[string]*User)
	storage.mTokens = make(map[string]*Token)
	storage.mProviders = make(map[string]d.Provider)
	storage.mMutex = &sync.Mutex{}
	storage.mIdentityID = 1
	return &storage
}

func (s *Storage) CreateRoom() (d.Room, error){
	roomID := atomic.AddInt32(&s.mIdentityID, 1)
	newRoom := NewRoom(roomID)
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	s.mRooms[roomID] = newRoom
	return newRoom, nil
}

func (s *Storage) GetRooms() ([]d.Room){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.Room, 0, len(s.mRooms))
	for _, v := range s.mRooms {
		ret = append(ret, v)
	}
	return ret
}

func (s *Storage) FindRoom(roomID int32) (d.Room, error){
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

func (s *Storage) DelRoom(roomID int32) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	_, ok := s.mRooms[roomID]
	if(!ok){
		return errors.New("No such room")
	}
	delete(s.mRooms, roomID)
	return nil
}

func (s *Storage) CreateUser(userName string, passwd string) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	_, ok := s.mUsers[userName]
	if(ok){
		return errors.New("UserName is occupied:" + userName)
	}
	s.mUsers[userName] = NewUser(userName, passwd)
	return nil
}

func (s *Storage) GetUsers() ([]d.User){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.User, 0, len(s.mUsers))
	for _, v := range s.mUsers {
		ret = append(ret, v)
	}
	return ret
}

func (s *Storage) FindUser(userName string, passwd string) (d.User, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	user, ok := s.mUsers[userName]
	if !ok {
		return nil, errors.New("No such user:" + userName)
	}
	return user, nil
}

func (s *Storage) DelUser(userName string) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	// TODO: is user playing?
	delete(s.mUsers, userName)
	return nil
}

func (s *Storage) CreateToken(user d.User) (d.Token, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	for k, token := range s.mTokens {
		if token.GetUserName() == user.GetUserName() {
			delete(s.mTokens, k)
		}
	}
	uuid := uuid.New()
	newToken := &Token{mKey: uuid, mUser: user}
	s.mTokens[uuid] = newToken
	return newToken, nil
}

func (s *Storage) FindToken(key string) (d.Token, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	token, ok := s.mTokens[key]
	if !ok {
		return nil, errors.New("No such token:" + key)
	}
	return token, nil
}

func (s *Storage) DelToken(key string) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	delete(s.mTokens, key)
	return nil
}

func (s *Storage) RegisterProvider(provider d.Provider) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	s.mProviders[provider.GetID()] = provider
	return nil
}

func (s *Storage) GetProviders() ([]d.Provider){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.Provider, 0, len(s.mProviders))
	for _, v := range s.mProviders {
		ret = append(ret, v)
	}
	return ret
}

func (s *Storage) FindProvider(providerID string) (d.Provider, error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	p, ok := s.mProviders[providerID]
	if(!ok){
		return nil, errors.New("No such provider:" + providerID)
	}
	return p, nil
}

func (s *Storage) FindProviderByGame(gameName string) ([]d.Provider){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	ret := make([]d.Provider, 0, len(s.mProviders))
	for _, v := range s.mProviders {
		if(v.GetGameName() == gameName){
			ret = append(ret, v)
		}
	}
	return ret
}

func (s *Storage) UnRegisterProvider(p d.Provider) (error){
	s.mMutex.Lock()
	defer s.mMutex.Unlock()
	delete(s.mProviders, p.GetID())
	return nil
}