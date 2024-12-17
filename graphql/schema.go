package graphql

import (
	"context"
	"go-graphql-echo-crud/db"
	"go-graphql-echo-crud/models"

	"github.com/graphql-go/graphql"
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

var updateUserMutation = &graphql.Field{
	Type: userType,
	Args: graphql.FieldConfigArgument{
		"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		"name":  &graphql.ArgumentConfig{Type: graphql.String},
		"email": &graphql.ArgumentConfig{Type: graphql.String},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id := params.Args["id"].(int)
		name, nameOk := params.Args["name"].(string)
		email, emailOk := params.Args["email"].(string)

		if nameOk && emailOk {
			_, err := db.Pool.Exec(context.Background(), "UPDATE users SET name=$1, email=$2 WHERE id=$3", name, email, id)
			if err != nil {
				return nil, err
			}
		} else if nameOk {
			_, err := db.Pool.Exec(context.Background(), "UPDATE users SET name=$1 WHERE id=$2", name, id)
			if err != nil {
				return nil, err
			}
		} else if emailOk {
			_, err := db.Pool.Exec(context.Background(), "UPDATE users SET email=$1 WHERE id=$2", email, id)
			if err != nil {
				return nil, err
			}
		}

		return models.User{ID: id, Name: name, Email: email}, nil
	},
}

var deleteUserMutation = &graphql.Field{
	Type: graphql.Boolean,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id := params.Args["id"].(int)

		_, err := db.Pool.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id)
		if err != nil {
			return nil, err
		}

		return true, nil
	},
}

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
			"updateUser": updateUserMutation,
			"deleteUser": deleteUserMutation,
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: mutation,
	},
)
