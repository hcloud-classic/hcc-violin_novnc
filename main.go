package main

import (
	"log"

	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/driver/grpccli"
	"hcc/violin-novnc/driver/grpcsrv"
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
	grpccli.CleanGRPCClient()

}

func main() {
	defer end()

	grpccli.InitGRPCClient()
	driver.VNCM.Prepare() // need harp to create proxy
	grpcsrv.InitGRPCServer()
}
