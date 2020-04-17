package main

import (
	"hcc/violin-novnc/action/graphql"
	"hcc/violin-novnc/lib/logger"
)

func main() {

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
		return
	}
}

