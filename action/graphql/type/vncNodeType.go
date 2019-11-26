package graphqlType

import (
	"github.com/graphql-go/graphql"
)

// VncNodeType: Form
var VncNodeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Novnc",
		Fields: graphql.Fields{
			"server_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"target_ip": &graphql.Field{
				Type: graphql.String,
			},
			"target_port": &graphql.Field{
				Type: graphql.String,
			},
			"target_pass": &graphql.Field{
				Type: graphql.String,
			},
			"websocket_port": &graphql.Field{
				Type: graphql.String,
			},
			"ws_url": &graphql.Field{
				Type: graphql.String,
			},
			"vnc_info": &graphql.Field{
				Type: graphql.String,
			},
			"action": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
