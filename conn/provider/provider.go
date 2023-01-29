package provider

import (
	"fmt"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	cpb "github.com/SailGame/Core/pb/core"
)

type Handler interface {
	HandleRegisterArgs(*Conn, *cpb.ProviderMsg, *cpb.RegisterArgs) error
	HandleNotifyMsg(*Conn, *cpb.ProviderMsg, *cpb.NotifyMsgArgs) error
	HandleCloseGameMsg(*Conn, *cpb.ProviderMsg, *cpb.CloseGameArgs) error
	Disconnect(*Conn)
}

type Conn struct {
	// bind to an entity
	ID interface{}
	// used in logging
	PrintID string

	mRunning    atomic.Value
	mServer     cpb.GameCore_ProviderServer
	mHandler    Handler
}

func NewConn(pServer cpb.GameCore_ProviderServer, handler Handler) *Conn {
	conn := &Conn{
		mServer:  pServer,
		mHandler: handler,
	}
	conn.mRunning.Store(false)
	return conn
}

func (conn *Conn) RecvLoop() {
	for {
		msg, err := conn.mServer.Recv()
		if err != nil {
			log.Warnf("Provider (%s) disconnected (%s)", conn.PrintID, err.Error())
			conn.mHandler.Disconnect(conn)
			break
		}
		log.Debugf("Provider connection (%s) recv msg (%s)", conn.PrintID, msg.String())

		if submsg := msg.GetRegisterArgs(); submsg != nil {
			err = conn.mHandler.HandleRegisterArgs(conn, msg, submsg)
		} else if submsg := msg.GetNotifyMsgArgs(); submsg != nil {
			err = conn.mHandler.HandleNotifyMsg(conn, msg, submsg)
		} else if submsg := msg.GetCloseGameArgs(); submsg != nil {
			err = conn.mHandler.HandleCloseGameMsg(conn, msg, submsg)
		} else {
			log.Warnf("Received unwanted msg (%s) from provider (%s)", msg.String(), conn.PrintID)
			conn.mHandler.Disconnect(conn)
			break
		}
		if err != nil {
			log.Infof("Process msg(%s) from provider (%s) failed. Err(%s)", msg.String(), conn.PrintID, err.Error())
			conn.mHandler.Disconnect(conn)
			break
		}
	}
}

func (conn *Conn) Poll() error {
	if conn.mRunning.CompareAndSwap(false, true) {
		log.Debugf("Provider connection started")
		conn.RecvLoop()
		return nil
	}
	return fmt.Errorf("Provider connection (%s) is running already", conn.PrintID)
}

func (conn *Conn) Send(msg *cpb.ProviderMsg) error {
	log.Debugf("Provider connection (%s) send msg (%s)", conn.PrintID, msg.String())
	return conn.mServer.Send(msg)
}
