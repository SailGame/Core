package server

import (
	"errors"
	"fmt"

	"github.com/SailGame/Core/conn/client"
	"github.com/SailGame/Core/conn/provider"
	d "github.com/SailGame/Core/data"
	cpb "github.com/SailGame/Core/pb/core"
)

func (coreServer *CoreServer) Provider(pServer cpb.GameCore_ProviderServer) error {
	conn := provider.NewConn(pServer, coreServer)
	conn.Start()
	return nil
}

func (coreServer *CoreServer) HandleRegisterArgs(conn *provider.Conn, providerMsg *cpb.ProviderMsg, regArgs *cpb.RegisterArgs) error {
	p := d.NewCommonProvider(conn, regArgs.Id, regArgs.GameName)
	if err := coreServer.mStorage.RegisterProvider(p); err != nil {
		return err
	}
	conn.ID = p
	conn.PrintID = regArgs.Id + ":" + regArgs.GameName
	return nil
}

func (coreServer *CoreServer) HandleNotifyMsg(conn *provider.Conn, providerMsg *cpb.ProviderMsg, notifyMsg *cpb.NotifyMsg) error {
	if conn.ID == nil {
		// not registered, ignore and disconnect
		return errors.New("Provider hasn't registered but sent other msgs")
	}
	p := conn.ID.(d.Provider)
	room := p.GetRoom(notifyMsg.RoomId)
	if room == nil {
		return errors.New(fmt.Sprintf("NotifyMsg: Unknown RoomId: (%d) conn: (%s)", notifyMsg.RoomId, conn.PrintID))
	}

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
			return errors.New(fmt.Sprintf("NotifyMsg: User (%s) disconnected", user.GetUserName()))
		}
		clientConn := conn.(*client.Conn)
		if (notifyMsg.UserId == 0) || (notifyMsg.UserId > 0 && uint32(notifyMsg.UserId) == user.GetTemporaryID()) || (uint32(-notifyMsg.UserId) != user.GetTemporaryID()) {
			err = clientConn.Send(broadcastMsg)
			if err != nil {
				user.SetConn(nil)
				return errors.New(fmt.Sprintf("NotifyMsg: User (%s) disconnected", user.GetUserName()))
			}
		}
	}

	return nil
}

func (coreServer *CoreServer) Disconnect(conn *provider.Conn) {
	if conn.ID != nil {
		p := conn.ID.(d.Provider)
		coreServer.mStorage.UnRegisterProvider(p)
	}
}
