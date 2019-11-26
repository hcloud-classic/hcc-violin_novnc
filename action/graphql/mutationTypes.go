package graphql

import "github.com/graphql-go/graphql"

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{},
})