package server

import (
	"log"

	d "github.com/SailGame/Core/data"
	cpb "github.com/SailGame/Core/pb/core"
)
func toGrpcRooms(rooms []d.Room) ([]*cpb.Room){
	ret := make([]*cpb.Room, 0, len(rooms))
	for _, v := range rooms {
		ret = append(ret, toGrpcRoom(v))
	}
	return ret
}

func toGrpcRoom(room d.Room) (*cpb.Room){
	return &cpb.Room{
		GameName: room.GetGameName(),
		RoomId: room.GetRoomID(),
		UserName: toUserName(room.GetUsers()),
	}
}

func toGrpcRoomDetails(room d.Room) (*cpb.RoomDetails, error){
	grpcRoomUsers, err := toGrpcRoomUsers(room)
	if err != nil{
		return nil, err
	}
	return &cpb.RoomDetails{
		GameName: room.GetGameName(),
		RoomId: room.GetRoomID(),
		User: grpcRoomUsers,
	}, nil
}

func toGrpcRoomUsers(room d.Room) ([]*cpb.RoomUser, error){
	users := room.GetUsers()
	ret := make([]*cpb.RoomUser, 0, len(users))
	for _, user := range users {
		state, err := room.GetUserState(user)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &cpb.RoomUser{
			UserName: user.GetUserName(),
			UserState: toGrpcRoomUserState(state),
		})
	}
	return ret, nil
}

func toGrpcRoomUserState(state d.UserState) (cpb.RoomUser_RoomUserState){
	if state == d.UserState_PREPARING{
		return cpb.RoomUser_PREPARING
	}else if state == d.UserState_READY{
		return cpb.RoomUser_READY
	}else if state == d.UserState_PLAYING{
		return cpb.RoomUser_PLAYING
	}else if state == d.UserState_EXITED{
		return cpb.RoomUser_DISCONNECTED
	}
	return cpb.RoomUser_ERROR
}

func toUserName(users []d.User) ([]string){
	ret := make([]string, 0, len(users))
	for _, v := range users {
		grpcUser := v.GetUserName()
		ret = append(ret, grpcUser)
	}
	return ret
}

func toBool(ready cpb.Ready) (bool){
	if(ready == cpb.Ready_READY){
		return true
	}
	if(ready == cpb.Ready_CANCEL){
		return false
	}
	log.Fatalf("Can't convert cpb.Ready %d", ready)
	return false
}

func toUserTempID(users []d.User) ([]uint32){
	ret := make([]uint32, 0, len(users))
	tid := uint32(1)
	for _, user := range users {
		user.SetTemporaryID(tid)
		ret = append(ret, tid)
		tid = tid + 1
	}
	return ret
}