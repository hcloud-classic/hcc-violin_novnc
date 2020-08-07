package main

import (
	"log"
	//"sync"

	//"hcc/violin-novnc/action/graphql"
	//"hcc/violin-novnc/driver"
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
}

func main() {
	defer end()
	// driver.RunProcxy("/var/log/violin-novnc/recordings/a/", "172.18.0.1:5901", "qwe1212", "5905")

	//http.Handle("/graphql", graphql.GraphqlHandler)
	grpccli.InitGRPCClient()
	grpcsrv.InitGRPCServer()
}
