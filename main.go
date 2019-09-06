package main

import (
	"GraphQL_violin_novnc/violin_novnccheckroot"
	"GraphQL_violin_novnc/violin_novncconfig"
	"GraphQL_violin_novnc/violin_novncgraphql"
	"GraphQL_violin_novnc/violin_novnclogger"
	"GraphQL_violin_novnc/violin_novncmysql"
	"net/http"
)

func main() {
	if !violin_novnccheckroot.CheckRoot() {
		return
	}

	if !violin_novnclogger.Prepare() {
		return
	}
	defer violin_novnclogger.FpLog.Close()

	err := violin_novncmysql.Prepare()
	if err != nil {
		return
	}
	defer violin_novncmysql.Db.Close()

	http.Handle("/graphql", violin_novncgraphql.GraphqlHandler)

	violin_novnclogger.Logger.Println("Server is running on port " + violin_novncconfig.HTTPPort)
	err = http.ListenAndServe(":"+violin_novncconfig.HTTPPort, nil)
	if err != nil {
		violin_novnclogger.Logger.Println("Failed to prepare http server!")
	}
}
