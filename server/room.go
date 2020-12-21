package server

import (
	"context"
	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer CoreServer) CreateRoom(ctx context.Context, req *cpb.CreateRoomArgs) (*cpb.CreateRoomRet, error) {
	_, err := coreServer.mStorage.CreateRoom()
	if(err != nil){
		return &cpb.CreateRoomRet{Errno: cpb.ErrorNumber_UnkownError}, nil
	}
	return &cpb.CreateRoomRet{Errno: cpb.ErrorNumber_OK}, nil
}

func (coreServer CoreServer) ControlRoom(ctx context.Context, req *cpb.ControlRoomArgs) (*cpb.ControlRoomRet, error) {
	room, err := coreServer.mStorage.FindRoom(req.RoomId)
	if(err != nil){
		return &cpb.ControlRoomRet{Errno: cpb.ErrorNumber_ControlRoom_RoomNotExist}, nil
	}
	if(req.GameName != ""){
		providers := coreServer.mStorage.FindProviderByGame(req.GameName)
		if(len(providers) == 0){
			return &cpb.ControlRoomRet{Errno: cpb.ErrorNumber_ControlRoom_RequiredProviderNotExist}, nil
		}
		// TODO: provider selector
		room.SetProvider(providers[0])
	}

	// TODO: passwd
	return &cpb.ControlRoomRet{Errno: cpb.ErrorNumber_OK}, nil
}

func (coreServer CoreServer) ListRoom(ctx context.Context, req *cpb.ListRoomArgs) (*cpb.ListRoomRet, error) {
	// TODO: game name filter
	return &cpb.ListRoomRet{Errno: cpb.ErrorNumber_OK, Room: toGrpc(coreServer.mStorage.GetRooms())}, nil
}

func (coreServer CoreServer) JoinRoom(ctx context.Context, req *cpb.JoinRoomArgs) (*cpb.JoinRoomRet, error) {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if(err != nil){
		return &cpb.JoinRoomRet{Errno: cpb.ErrorNumber_JoinRoom_InvalidToken}, nil
	}
	room, err := coreServer.mStorage.FindRoom(req.RoomId)
	if(err != nil){
		return &cpb.JoinRoomRet{Errno: cpb.ErrorNumber_JoinRoom_InvalidRoomID}, nil
	}
	err = room.UserJoin(token.GetUser())
	if(err != nil){
		// TODO: clearer error?
		return &cpb.JoinRoomRet{Errno: cpb.ErrorNumber_JoinRoom_FullRoom}, nil
	}
	return &cpb.JoinRoomRet{Errno: cpb.ErrorNumber_OK}, nil
}

func (coreServer CoreServer) ExitRoom(ctx context.Context, req *cpb.ExitRoomArgs) (*cpb.ExitRoomRet, error) {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if(err != nil){
		return &cpb.ExitRoomRet{Errno: cpb.ErrorNumber_JoinRoom_InvalidToken}, nil
	}
	token.GetUser().SetRoom(nil)
	return &cpb.ExitRoomRet{Errno: cpb.ErrorNumber_OK}, nil
}

// func (coreServer CoreServer) QueryRoom(ctx context.Context, req *cpb.QueryRoomArgs) (*cpb.QueryRoomRet, error) {
// 	return nil, nil
// }

// func (coreServer CoreServer) OperationInRoom(ctx context.Context, req *cpb.OperationInRoomArgs) (*cpb.OperationInRoomRet, error) {
// 	return nil, nil
// }
