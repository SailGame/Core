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