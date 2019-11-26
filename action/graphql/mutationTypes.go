package graphql

import (
	"github.com/graphql-go/graphql"
	graphqlType "hcc/violin-novnc/action/graphql/type"
	"hcc/violin-novnc/driver"
	"hcc/violin-novnc/lib/logger"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create_vnc": &graphql.Field{
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
				logger.Logger.Println("Resolving: create_vnc")
				return driver.Runner(params)
			},
		},
	},
})