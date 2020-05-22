package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Schema : GraphQL schema definition
var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryTypes,
		Mutation: mutationTypes,
	},
)

// GraphqlHandler : Show GraphQL GUI request form in web browser
var GraphqlHandler = handler.New(&handler.Config{
	Schema:   &Schema,
	Pretty:   true,
	GraphiQL: true,
})
