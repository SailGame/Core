package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/SailGame/Core/conn/client"
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
	log.Infof("Control Room token(%s) GameName(%s)", req.GetToken(), req.GetGameName())

	room, err := coreServer.mStorage.FindRoom(req.RoomId)
	if err != nil {
		return &cpb.ControlRoomRet{Err: cpb.ErrorNumber_ControlRoom_RoomNotExist}, nil
	}
	room.Lock()
	defer room.Unlock()
	if req.GameName != "" {
		providers := coreServer.mStorage.FindProviderByGame(req.GameName)
		if len(providers) == 0 {
			return &cpb.ControlRoomRet{Err: cpb.ErrorNumber_ControlRoom_RequiredProviderNotExist}, nil
		}
		// TODO: provider selector
		room.SetProvider(providers[0])
		room.SetCustomGameSetting(req.GetCustom())
	}

	coreServer.NotifyRoomDetails(ctx, room)
	// TODO: passwd
	return &cpb.ControlRoomRet{Err: cpb.ErrorNumber_OK}, nil
}

func (coreServer *CoreServer) ListRoom(ctx context.Context, req *cpb.ListRoomArgs) (*cpb.ListRoomRet, error) {
	log.Debugf("List Room gameName(%s)", req.GetGameName())
	filter := func(r d.Room) bool {
		return req.GetGameName() == "" || r.GetGameName() == req.GetGameName()
	}
	return &cpb.ListRoomRet{Err: cpb.ErrorNumber_OK, Room: toGrpcRooms(coreServer.mStorage.GetRoomsWithFilter(filter))}, nil
}

func (coreServer *CoreServer) JoinRoom(ctx context.Context, req *cpb.JoinRoomArgs) (*cpb.JoinRoomRet, error) {
	log.Debugf("Join Room token(%s) roomID(%d)", req.GetToken(), req.GetRoomId())

	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_InvalidToken}, nil
	}
	user := token.GetUser()
	user.Lock()
	defer user.Unlock()

	curRoom, err := user.GetRoom()
	if curRoom != nil {
		if curRoom.GetRoomID() == req.RoomId {
			return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_OK}, nil
		} else {
			return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_UserIsInAnotherRoom}, nil
		}
	}

	room, err := coreServer.mStorage.FindRoom(req.RoomId)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_InvalidRoomID}, nil
	}
	room.Lock()
	defer room.Unlock()

	err = room.UserJoin(user)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_JoinRoom_FullRoom}, nil
	}
	err = user.SetRoom(room)
	if err != nil {
		return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	log.Debugf("Join Room userName(%s) roomID(%d)", user.GetUserName(), req.GetRoomId())

	coreServer.NotifyRoomDetails(ctx, room)

	return &cpb.JoinRoomRet{Err: cpb.ErrorNumber_OK}, nil
}

func (coreServer *CoreServer) ExitRoom(ctx context.Context, req *cpb.ExitRoomArgs) (*cpb.ExitRoomRet, error) {
	log.Infof("Exit Room token(%s)", req.GetToken())

	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_ExitRoom_InvalidToken}, nil
	}
	user := token.GetUser()
	user.Lock()
	defer user.Unlock()
	room, err := user.GetRoom()
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_ExitRoom_NotInRoom}, nil
	}
	room.Lock()
	defer room.Unlock()
	ok, err := room.UserExit(user)
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	if !ok {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_ExitRoom_IsPlaying}, nil
	}
	err = user.SetRoom(nil)
	if err != nil {
		return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	log.Infof("Exit Room userName(%s) roomID(%d)", user.GetUserName(), room.GetRoomID())

	coreServer.NotifyRoomDetails(ctx, room)

	return &cpb.ExitRoomRet{Err: cpb.ErrorNumber_OK}, nil
}

func (coreServer *CoreServer) QueryRoom(ctx context.Context, req *cpb.QueryRoomArgs) (*cpb.QueryRoomRet, error) {
	room, err := coreServer.mStorage.FindRoom(req.GetRoomId())
	if err != nil {
		return &cpb.QueryRoomRet{Err: cpb.ErrorNumber_QryRoom_InvalidRoomID}, nil
	}
	room.Lock()
	defer room.Unlock()
	roomDetails, err := toGrpcRoomDetails(room)
	if err != nil {
		return &cpb.QueryRoomRet{Err: cpb.ErrorNumber_UnkownError}, nil
	}
	return &cpb.QueryRoomRet{Err: cpb.ErrorNumber_OK, Room: roomDetails}, nil
}

func (coreServer *CoreServer) OperationInRoom(ctx context.Context, req *cpb.OperationInRoomArgs) (*cpb.OperationInRoomRet, error) {
	token, err := coreServer.mStorage.FindToken(req.Token)
	if err != nil {
		return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OperRoom_InvalidToken}, nil
	}
	token.GetUser().Lock()
	defer token.GetUser().Unlock()
	room, err := token.GetUser().GetRoom()
	if err != nil {
		return &cpb.OperationInRoomRet{Err: cpb.ErrorNumber_OperRoom_NotInRoom}, nil
	}
	room.Lock()
	defer room.Unlock()
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
		coreServer.NotifyRoomDetails(ctx, room)
		if room.GetState() == d.RoomState_PLAYING {
			pConn.Send(&cpb.ProviderMsg{
				Msg: &cpb.ProviderMsg_StartGameArgs{
					StartGameArgs: &cpb.StartGameArgs{
						RoomId: room.GetRoomID(),
						UserId: toUserTempID(room.GetUsers()),
						Custom: room.GetCustomGameSetting(),
					},
				},
			})
		}
	} else if custom := req.GetCustom(); custom != nil {
		pConn.Send(&cpb.ProviderMsg{
			Msg: &cpb.ProviderMsg_UserOperationArgs{
				UserOperationArgs: &cpb.UserOperationArgs{
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

func (coreServer *CoreServer) NotifyRoomDetails(ctx context.Context, room d.Room) error {
	roomDetails, err := toGrpcRoomDetails(room)
	if err != nil {
		return err
	}
	for _, v := range room.GetUsers() {
		conn, err := v.GetConn()
		if err != nil {
			return err
		}
		clientConn := conn.(*client.Conn)
		err = clientConn.Send(&cpb.BroadcastMsg{
			Msg: &cpb.BroadcastMsg_RoomDetails{
				RoomDetails: roomDetails,
			},
		})
		if err != nil {
			return nil
		}
	}
	return nil
}
