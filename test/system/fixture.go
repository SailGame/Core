package system

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	cpb "github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/server"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type systemFixture struct {
	mCoreServer *server.CoreServer
	mLis        *bufconn.Listener
	mCtx        context.Context
	mConn       *grpc.ClientConn
}

type userClient struct {
	mCoreClient cpb.GameCoreClient
	mLisClient  cpb.GameCore_ListenClient
	mClose      context.CancelFunc
}

func (uc *userClient) listenToCore(token string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	uc.mClose = cancel
	uc.mLisClient, err = uc.mCoreClient.Listen(ctx, &cpb.ListenArgs{
		Token: token,
	})
	return
}

type providerClient struct {
	mCoreClient     cpb.GameCoreClient
	mProviderClient cpb.GameCore_ProviderClient
}

func (pc *providerClient) connectToCore() (err error) {
	pc.mProviderClient, err = pc.mCoreClient.Provider(context.Background())
	return
}

func newFixture() *systemFixture {
	const bufSize = 1024 * 1024
	return &systemFixture{
		mCtx: context.Background(),
		mLis: bufconn.Listen(bufSize),
	}
}

func (sf *systemFixture) init(config *server.CoreServerConfig) error {
	s := grpc.NewServer()
	var err error
	sf.mCoreServer, err = server.NewCoreServer(config)
	if err != nil {
		return err
	}
	cpb.RegisterGameCoreServer(s, sf.mCoreServer)
	go func() {
		if err := s.Serve(sf.mLis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	return nil
}

func (sf *systemFixture) dial() *grpc.ClientConn {
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return sf.mLis.Dial()
	}
	conn, err := grpc.DialContext(sf.mCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}

func (sf *systemFixture) newUserClient() *userClient {
	conn := sf.dial()
	return &userClient{
		mCoreClient: cpb.NewGameCoreClient(conn),
	}
}

func (sf *systemFixture) newProviderClient() *providerClient {
	conn := sf.dial()
	return &providerClient{
		mCoreClient: cpb.NewGameCoreClient(conn),
	}
}
