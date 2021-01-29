package main

import (
	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/action/grpc/server"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/mysql"
)

func init() {
	err := logger.Init()
	if err != nil {
		err.Fatal()
	}

	config.Parser()

	err = mysql.Init()
	if err != nil {
		err.Fatal()
	}

	client.InitGRPCClient()

	es := driver.Init()
	if es != nil {
		logger.Logger.Println("noVNCDriver Init failed. Skip proxy restore")
		es.Dump()
	}
}

func end() {
	logger.End()
	mysql.End()
	client.CleanGRPCClient()
	server.CleanGRPCServer()
}

func main() {
	defer end()

	server.InitGRPCServer()
}
