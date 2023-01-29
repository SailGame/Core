package server

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/SailGame/Core/conn/client"
	"github.com/SailGame/Core/conn/provider"
	d "github.com/SailGame/Core/data"
	cpb "github.com/SailGame/Core/pb/core"
)

// settings
func getMaxUsers(regArgs *cpb.RegisterArgs) int32 {
	if regArgs == nil || regArgs.GameSetting == nil || regArgs.GameSetting.MaxUsers <= 0 {
		return 0
	}
	return regArgs.GameSetting.MaxUsers
}

func getMinUsers(regArgs *cpb.RegisterArgs) int32 {
	if regArgs == nil || regArgs.GameSetting == nil || regArgs.GameSetting.MaxUsers <= 0 {
		return 0
	}
	return regArgs.GameSetting.MinUsers
}

func (coreServer *CoreServer) Provider(pServer cpb.GameCore_ProviderServer) error {
	log.Info("Provider connected")
	conn := provider.NewConn(pServer, coreServer)
	conn.Poll()
	return nil
}

func (coreServer *CoreServer) HandleRegisterArgs(conn *provider.Conn, providerMsg *cpb.ProviderMsg, regArgs *cpb.RegisterArgs) error {
	p := d.NewCommonProvider(conn, regArgs.GetId(), d.CommonGameSetting{
		GameName: regArgs.GetGameName(),
		MaxUser:  getMaxUsers(regArgs),
		MinUser:  getMinUsers(regArgs),
	})
	if err := coreServer.mStorage.RegisterProvider(p); err != nil {
		return fmt.Errorf("Provider register failed: (%s) (%s) (%s)", regArgs.GetId(), regArgs.GetGameName(), err.Error())
	}
	log.Infof("Provider register successfully: (%s) (%s)", regArgs.GetId(), regArgs.GetGameName())
	conn.ID = p
	conn.PrintID = regArgs.GetId() + ":" + regArgs.GetGameName()

	conn.Send(&cpb.ProviderMsg{
		Msg: &cpb.ProviderMsg_RegisterRet{
			RegisterRet: &cpb.RegisterRet{
				Err: cpb.ErrorNumber_OK,
			},
		},
	})
	return nil
}

func (coreServer *CoreServer) HandleNotifyMsg(conn *provider.Conn, providerMsg *cpb.ProviderMsg, notifyMsg *cpb.NotifyMsgArgs) error {
	if conn.ID == nil {
		// not registered, ignore and disconnect
		return errors.New("Provider hasn't registered but sent other msgs")
	}

	p := conn.ID.(d.Provider)
	room := p.GetRoom(notifyMsg.RoomId)
	if room == nil {
		return fmt.Errorf("NotifyMsgArgs: Unknown RoomId: (%d) conn: (%s)", notifyMsg.RoomId, conn.PrintID)
	}

	log.Debugf("NotifyMsgArgs: conn(%s) room(%d)", conn.PrintID, notifyMsg.RoomId)

	broadcastMsg := &cpb.BroadcastMsg{
		FromUser: 0,
		ToUser:   0,
		Msg: &cpb.BroadcastMsg_Custom{
			Custom: notifyMsg.Custom,
		},
	}
	for _, user := range room.GetUsers() {
		conn, err := user.GetConn()
		if err != nil {
			log.Warnf("User(%s) error: %v", user.GetUserName(), err)
			return nil
		}
		clientConn := conn.(*client.Conn)
		if (notifyMsg.UserId == 0) || (notifyMsg.UserId > 0 && uint32(notifyMsg.UserId) == user.GetTemporaryID()) || (notifyMsg.UserId < 0 && uint32(-notifyMsg.UserId) != user.GetTemporaryID()) {
			err = clientConn.Send(broadcastMsg)
			if err != nil {
				user.SetConn(nil)
				log.Warnf("Client(%s) Disconnect", user.GetUserName())
				return nil
			}
		}
	}

	return nil
}

func (coreServer *CoreServer) HandleCloseGameMsg(conn *provider.Conn, providerMsg *cpb.ProviderMsg, closeGameArgs *cpb.CloseGameArgs) error {
	if conn.ID == nil {
		// not registered, ignore and disconnect
		return errors.New("Provider hasn't registered but sent other msgs")
	}
	roomID := closeGameArgs.GetRoomId()
	room, err := coreServer.mStorage.FindRoom(roomID)
	if err != nil {
		return fmt.Errorf("CloseGameArgs: Unknown RoomId: (%d) conn: (%s)", roomID, conn.PrintID)
	}
	log.Debugf("CloseGameMsg: conn(%s) room(%d)", conn.PrintID, closeGameArgs.RoomId)

	room.Lock()
	defer room.Unlock()
	room.Restart()
	return nil
}

func (coreServer *CoreServer) Disconnect(conn *provider.Conn) {
	if conn.ID != nil {
		p := conn.ID.(d.Provider)
		coreServer.mStorage.UnRegisterProvider(p)
	}
}
