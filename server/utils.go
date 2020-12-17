package server

import (
	cpb "github.com/SailGame/Core/pb/core"
	d "github.com/SailGame/Core/data"
)
func toGrpc(rooms []d.Room) ([]*cpb.Room){
	ret := make([]*cpb.Room, len(rooms))
	for _, v := range rooms {
		grpcRoom := cpb.Room{
			GameName: v.GetGameName(),
			RoomId: v.GetRoomID(),
			UserName: toUserName(v.GetUsers()),
		}
		ret = append(ret, &grpcRoom)
	}
	return ret
}

func toUserName(users []d.User) ([]string){
	ret := make([]string, len(users))
	for _, v := range users {
		grpcUser := v.GetUserName()
		ret = append(ret, grpcUser)
	}
	return ret
}