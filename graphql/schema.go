package graphql

import (
	"context"
	"github.com/graphql-go/graphql"
	"go-graphql-echo-crud/db"
	"go-graphql-echo-crud/models"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},
			"name":  &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
		},
	},
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					rows, err := db.Pool.Query(context.Background(), "SELECT id, name, email FROM users")
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					var users []models.User
					for rows.Next() {
						var user models.User
						err := rows.Scan(&user.ID, &user.Name, &user.Email)
						if err != nil {
							return nil, err
						}
						users = append(users, user)
					}
					return users, nil
				},
			},
		},
	},
)

var mutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"name":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name := params.Args["name"].(string)
					email := params.Args["email"].(string)

					_, err := db.Pool.Exec(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
					if err != nil {
						return nil, err
					}

					return models.User{Name: name, Email: email}, nil
				},
			},
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: mutation,
	},
)
