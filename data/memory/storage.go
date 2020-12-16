package memory

import (
	d "github.com/SailGame/Core/data"
)

type Storage struct {
	map[int]Room mRooms
	map[string]User mUsers
	map[string]Token mTokens
	map[string]interface{} mProviders
}

func NewStorage() (Storage)
{
	storage := Storage{}
	storage
}

func (s *Storage) CreateRoom(roomID int) (error)
{
	
}
func (s *Storage) FindRoom(roomID int) (d.Room, error)
{

}
func (s *Storage) DelRoom(roomID int) (error)
{

}

func (s *Storage) CreateUser(userName string, passwd string) (error)
{

}
func (s *Storage) FindUser(userName string, passwd string) (d.User, error)
{

}
func (s *Storage) DelUser(userName string) (error)
{

}

func (s *Storage) CreateToken(userID int) (error)
{

}
func (s *Storage) FindToken(token string) (d.Token, error)
{

}
func (s *Storage) DelToken(token string) (error)
{

}
func (s *Storage) RegisterProvider(providerID string, provider interface{}) (error)
{

}
func (s *Storage) FindProvider(providerID string) (interface{}, error)
{

}
func (s *Storage) UnRegisterProvider(providerID string) (error)
{

}

type Room struct {

}

type User struct {

}

type Token struct {

}