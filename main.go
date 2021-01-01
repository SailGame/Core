package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/SailGame/Core/data/memory"
	cpb "github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	serverAddr = flag.String("serverAddr", "0.0.0.0:8080", "The server address in the format of host:port")
	logFile    = flag.String("logFile", "", "The log file of GoDock")
	logLevel   = flag.String("logLevel", "Info", "The log level. (Debug/Info/Warn/Error)")
)

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	if *logFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("error opening logFile(%s): %v", *logFile, err))
		}
		log.SetOutput(f)
	}

	if *logLevel == "Debug" {
		log.SetLevel(log.DebugLevel)
	} else if *logLevel == "Info" {
		log.SetLevel(log.InfoLevel)
	} else if *logLevel == "Warn" {
		log.SetLevel(log.WarnLevel)
	} else if *logLevel == "Error" {
		log.SetLevel(log.ErrorLevel)
	} else {
		panic(fmt.Sprintf("Unknown logLevel: %s", *logLevel))
	}
}

func main() {
	flag.Parse()
	initLog()
	lis, err := net.Listen("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	coreServer, err := server.NewCoreServer(&server.CoreServerConfig{
		MStorage: memory.NewStorage(),
	})
	if err != nil {
		log.Fatalf("failed to create core server: %v", err)
	}
	cpb.RegisterGameCoreServer(s, coreServer)
	reflection.Register(s)
	log.Infof("Start Game Core Server at %s", *serverAddr)
	log.Fatal(s.Serve(lis))
}
