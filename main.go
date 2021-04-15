package main

import (
	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/action/grpc/server"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
)

func init() {
	err := logger.Init()
	if err != nil {
		err.Fatal()
	}

	config.Parser()

}

func end() {
	logger.End()
	client.CleanGRPCClient()
	server.CleanGRPCServer()
}

func main() {
	defer end()

	client.InitGRPCClient()
	server.InitGRPCServer()
}
