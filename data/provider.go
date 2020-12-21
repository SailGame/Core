package data

import "sync"

// Provider supports at least one kind of game, and connects to core server
// as there is a live connection, provider is always in-memory
// right now we just use CommonProvider.
// however, should we bind special methods to different provider?
// e.g. UnoProvider, TexasProvider are both derived from CommonProvider
// and they have some custom functions.
type Provider interface {
	GetID() (string)
	GetGameName() (string)
	GetRooms() ([]Room)
	GetRoom(int32) (Room)

	AddRoom(Room) (error)
	DelRoom(Room) (error)
}

// CommonProvider has the basic functionality of a 'provider'
type CommonProvider struct {
	mID string
	mGameName string
	mRooms map[int32]Room
	mMutex sync.Locker
}

func NewCommonProvider(id string, gameName string) (*CommonProvider) {
	provider := &CommonProvider{
		mID: id,
		mGameName: gameName,
		mRooms: make(map[int32]Room),
		mMutex: &sync.Mutex{},
	}
	return provider
}

func (cp CommonProvider) GetID() (string){
	return cp.mID
}

func (cp CommonProvider) GetGameName() (string){
	return cp.mGameName
}

func (cp CommonProvider) GetRooms() ([]Room){
	ret := make([]Room, len(cp.mRooms))
	cp.mMutex.Lock()
	defer cp.mMutex.Unlock()
	for _, v := range cp.mRooms {
		ret = append(ret, v)
	}
	return ret
}

func (cp CommonProvider) GetRoom(roomId int32) (Room){
	room, ok := cp.mRooms[roomId]
	if(!ok){
		return nil
	}
	return room
}

func (cp CommonProvider) AddRoom(r Room) (error){
	cp.mMutex.Lock()
	defer cp.mMutex.Unlock()
	cp.mRooms[r.GetRoomID()] = r
	return nil
}

func (cp CommonProvider) DelRoom(r Room) (error){
	cp.mMutex.Lock()
	defer cp.mMutex.Unlock()
	// TODO: check existence
	delete(cp.mRooms, r.GetRoomID())
	return nil
}