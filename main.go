package main

import (
	"fmt"
	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/action/grpc/server"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
	"hcc/violin-novnc/lib/mysql"
	"hcc/violin-novnc/lib/pid"
	"os"
	"strconv"

	"innogrid.com/hcloud-classic/hcc_errors"
)

func init() {
	violinNovncRunning, violinNovncPID, err := pid.IsViolinNovncRunning()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if violinNovncRunning {
		fmt.Println("violin-novnc is already running. (PID: " + strconv.Itoa(violinNovncPID) + ")")
		os.Exit(1)
	}
	err = pid.WriteViolinNovncPID()
	if err != nil {
		_ = pid.DeleteViolinNovncPID()
		fmt.Println(err)
		panic(err)
	}

	err = logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.ViolinNoVNCInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteViolinNovncPID()
	}
	hcc_errors.SetErrLogger(logger.Logger)

	config.Init()

	err = client.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.ViolinNoVNCInternalInitFail, "client.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteViolinNovncPID()
	}

	err = mysql.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.ViolinNoVNCInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
		_ = pid.DeleteViolinNovncPID()
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
	_ = pid.DeleteViolinNovncPID()
}

func main() {
	defer end()

	server.InitGRPCServer()
}
