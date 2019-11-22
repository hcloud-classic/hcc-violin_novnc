package main

import (
	"fmt"
	"net/http"
	"strconv"

	"hcc/violin-novnc/action/graphql"
	vncEnd "hcc/violin-novnc/end"
	vncInit "hcc/violin-novnc/init"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
)

func init() {
	err := vncInit.MainInit()
	if err != nil {
		panic(err)
	}
}
func main() {
	defer func() {
		vncEnd.MainEnd()
	}()

	// driver.RunProcxy("/var/log/violin-novnc/recordings/a/", "172.18.0.1:5901", "qwe1212", "5905")
	fmt.Println(config.HTTP.Port)

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
		return
	}
}
