package provider

import (
	"sync"
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

	mRecvLoopWg sync.WaitGroup
	mRunning    atomic.Value
	mServer     cpb.GameCore_ProviderServer
	mHandler    Handler
	mMutex      sync.Locker
}

func NewConn(pServer cpb.GameCore_ProviderServer, handler Handler) *Conn {
	conn := &Conn{
		mServer:  pServer,
		mHandler: handler,
		mMutex:   &sync.Mutex{},
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
			conn.mRunning.Store(false)
			return
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
			conn.mRunning.Store(false)
			return
		}
		if err != nil {
			log.Infof("Process msg(%s) from provider (%s) fail (%s)", msg.String(), conn.PrintID, err.Error())
			conn.mHandler.Disconnect(conn)
			conn.mRunning.Store(false)
			return
		}
	}

}

func (conn *Conn) Start() {
	if conn.mRunning.Load().(bool) {
		return
	}
	log.Debugf("Provider connection started")
	conn.mRunning.Store(true)
	conn.mRecvLoopWg.Add(1)
	go conn.RecvLoop()
	conn.mRecvLoopWg.Wait()
}

func (conn *Conn) Close() {
	if !conn.mRunning.Load().(bool) {
		return
	}
	log.Debugf("Provider connection (%s) received stop signal ", conn.PrintID)
	conn.mRecvLoopWg.Done()
}

func (conn *Conn) Send(msg *cpb.ProviderMsg) error {
	log.Debugf("Provider connection (%s) send msg (%s)", conn.PrintID, msg.String())
	conn.mMutex.Lock()
	err := conn.mServer.Send(msg)
	conn.mMutex.Unlock()
	if err != nil {
		log.Debug(err)
	}
	return err
}
