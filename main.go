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
		defer logger.End()
		err.Fatal()
	}
}

func end() {
	mysql.End()
	logger.End()
	client.CleanGRPCClient()
	server.CleanGRPCServer()
}

func main() {
	defer end()

	client.InitGRPCClient()
	driver.VNCD.Prepare() // need harp to create proxy
	server.InitGRPCServer()
}
