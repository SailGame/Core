package main

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	cpb "github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/server"
	"google.golang.org/grpc"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	// TODO: cmd args
	lis, err := net.Listen("tcp", "0.0.0.0" + ":" + "8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	coreServer, err := server.NewCoreServer(&server.CoreServerConfig{})
	if err != nil {
		panic(err)
	}
	cpb.RegisterGameCoreServer(s, coreServer)
	log.Info("rpc server start")
	s.Serve(lis)
}