package graphql

import (
	graphqlType "hcc/violin-novnc/action/graphql/type"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/logger"

	"github.com/graphql-go/graphql"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"control_vnc": &graphql.Field{
			Type:        graphqlType.VncNodeType,
			Description: "Create vnc",
			Args: graphql.FieldConfigArgument{
				"server_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"target_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"target_port": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"target_pass": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"action": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: control_vnc")
				return driver.Runner(params)
			},
		},
	},
})
