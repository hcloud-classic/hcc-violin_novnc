package main

import (
	"log"

	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/action/grpc/server"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/mysql"
	"hcc/violin-novnc/lib/syscheck"
)

func init() {
	err := syscheck.CheckRoot()
	if err != nil {
		log.Panic(err)
	}

	if !logger.Prepare() {
		log.Panic("error occurred while preparing logger")
	}

	config.Parser()

	err = mysql.Prepare()
	if err != nil {
		logger.FpLog.Close()
		log.Panic(err)
	}
}

func end() {
	mysql.Db.Close()
	logger.FpLog.Close()
	client.CleanGRPCClient()
	server.CleanGRPCServer()

}

func main() {
	defer end()

	client.InitGRPCClient()
	driver.VNCD.Prepare() // need harp to create proxy
	server.InitGRPCServer()
}
