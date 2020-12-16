package main

import (
	"log"
	"net"

	"github.com/SailGame/Core/server"
	cpb "github.com/SailGame/Core/pb/core"
	"google.golang.org/grpc"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
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
	log.Println("rpc server start")
	s.Serve(lis)
}