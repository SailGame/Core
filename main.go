package main

import (
	"log"
	"net"

	"github.com/SailGame/Core/handler"
	cpb "github.com/SailGame/Core/pb/core"
	"google.golang.org/grpc"
)

func logInit() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	logInit()

	// TODO: cmd args
	lis, err := net.Listen("tcp", "0.0.0.0" + ":" + "8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	handler, err := handler.NewCoreServerHandler(&handler.CoreServerHandlerConfig{})
	if err != nil {
		panic(err)
	}
	cpb.RegisterGameCoreServer(s, handler)
	log.Println("rpc server start")
	s.Serve(lis)
}