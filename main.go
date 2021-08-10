package main

import (
	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/action/grpc/server"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/mysql"

	"innogrid.com/hcloud-classic/hcc_errors"
)

func init() {
	err := logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.ViolinNoVNCInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(logger.Logger)

	config.Init()

	err = mysql.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.ViolinNoVNCInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
	}

	err = client.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.ViolinNoVNCInternalInitFail, "client.Init(): "+err.Error()).Fatal()
	}

	es := driver.Init()
	if es != nil {
		logger.Logger.Println("noVNCDriver Init failed. Skip proxy restore")
		_ = es.Dump()
	}
}

func end() {
	server.CleanGRPCServer()
	client.End()
	mysql.End()
	logger.End()
}

func main() {
	defer end()

	server.InitGRPCServer()
}
