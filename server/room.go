package server

import (
	"context"

	"github.com/SailGame/Core/conn/provider"
	d "github.com/SailGame/Core/data"
	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer *CoreServer) CreateRoom(ctx context.Context, req *cpb.CreateRoomArgs) (*cpb.CreateRoomRet, error) {
	room, err := coreServer.mStorage.CreateRoom()
	if err != nil {
		return &cpb.CreateRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	return &cpb.CreateRoomRet{Err: cpb.ErrorNumber_OK, RoomId: room.GetRoomID()}, nil
}

func (coreServer *CoreServer) ControlRoom(ctx context.Context, req *cpb.ControlRoomArgs) (*cpb.ControlRoomRet, error) {
	room, err := coreServer.mStorage.FindRoom(req.RoomId)
	if err != nil {
		return &cpb.ControlRoomRet{Err: cpb.ErrorNumber_ControlRoom_RoomNotExist}, nil
	}
	if req.GameName != "" {
		providers := coreServer.mStorage.FindProviderByGame(req.GameName)
		if len(providers) == 0 {
			return &cpb.ControlRoomRet{Err: cpb.ErrorNumber_ControlRoom_RequiredProviderNotExist}, nil
		}
		// TODO: provider selector
		room.SetProvider(providers[0])
	}

	// TODO: passwd
	return &cpb.ControlRoomRet{Err: cpb.ErrorNumber_OK}, nil
}

func (coreServer *CoreServer) ListRoom(ctx context.Context, req *cpb.ListRoomArgs) (*cpb.ListRoomRet, error) {
	filter := func(r d.Room) bool {
		return req.GetGameName() == "" || r.GetGameName() == req.GetGameName()
	}
	return &cpb.ListRoomRet{Err: cpb.ErrorNumber_OK, Room: toGrpcRooms(coreServer.mStorage.GetRoomsWithFilter(filter))}, nil
}

func (coreServer *CoreServer) JoinRoom(ctx context.Context, req *cpb.JoinRoomArgs) (*cpb.JoinRoomRet, error) {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_InvalidToken}, nil
	}
	user := token.GetUser()
	room, err := coreServer.mStorage.FindRoom(req.RoomId)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_InvalidRoomID}, nil
	}
	curRoom, err := user.GetRoom()
	if curRoom == room {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_OK}, nil
	} else if curRoom != nil {
		err := curRoom.UserExit(user)
		if err != nil {
			return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_UserIsInAnotherRoomAndFailToExit}, nil
		}
	}
	err = room.UserJoin(user)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_FullRoom}, nil
	}
	err = user.SetRoom(room)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_OK}, nil
}

func (coreServer *CoreServer) ExitRoom(ctx context.Context, req *cpb.ExitRoomArgs) (*cpb.ExitRoomRet, error) {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_ExitRoom_InvalidToken}, nil
	}
	user := token.GetUser()
	room, err := user.GetRoom()
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_ExitRoom_NotInRoom}, nil
	}
	err = room.UserExit(user)
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	err = user.SetRoom(nil)
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_OK}, nil
}

func (coreServer *CoreServer) QueryRoom(ctx context.Context, req *cpb.QueryRoomArgs) (*cpb.QueryRoomRet, error) {
	room, err := coreServer.mStorage.FindRoom(req.GetRoomId())
	if err != nil {
		return &cpb.QueryRoomRet{Err: cpb.ErrorNumber_QryRoom_InvalidRoomID}, nil
	}
	return &cpb.QueryRoomRet{Err: cpb.ErrorNumber_OK, Room: toGrpcRoom(room)}, nil
}

func (coreServer *CoreServer) OperationInRoom(ctx context.Context, req *cpb.OperationInRoomArgs) (*cpb.OperationInRoomRet, error) {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OperRoom_InvalidToken}, nil
	}
	room, err := token.GetUser().GetRoom()
	if err != nil {
		return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OperRoom_NotInRoom}, nil
	}
	pv := room.GetProvider()
	if pv == nil {
		return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OperRoom_ProviderUnavailable}, nil
	}
	conn, err := pv.GetConn()
	if err != nil {
		return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OperRoom_ProviderUnavailable}, nil
	}
	pConn := conn.(*provider.Conn)

	if req.GetReady() != cpb.Ready_UNSET {
		err := room.UserReady(token.GetUser(), toBool(req.GetReady()))
		if err != nil {
			// room is in playing or ??
			return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OperRoom_CannotChangeReadyState}, nil
		}
		if room.GetState() == d.Playing {
			pConn.Send(&cpb.ProviderMsg{
				Msg: &cpb.ProviderMsg_StartGameArgs{
					&cpb.StartGameArgs{
						RoomId: room.GetRoomID(),
						UserId: toUserTempID(room.GetUsers()),
						Custom: nil,
					},
				},
			})
		}
	} else if custom := req.GetCustom(); custom != nil {
		pConn.Send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_UserOperationArgs{
				&cpb.UserOperationArgs{
					RoomId: room.GetRoomID(),
					UserId: token.GetUser().GetTemporaryID(),
					Custom: custom,
				},
			},
		})
	} else {
		return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}

	return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OK}, nil
}
