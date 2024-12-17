package main

import (
	"go-graphql-echo-crud/db"
	"go-graphql-echo-crud/graphql"

	"github.com/graphql-go/handler"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	// Initialize the database
	db.Init()

	// Initialize Echo
	e := echo.New()

	// GraphQL Handler
	graphqlHandler := handler.New(&handler.Config{
		Schema:   &graphql.Schema, // Correctly referenced Schema
		Pretty:   true,
		GraphiQL: true, // Enables GraphiQL UI for testing
	})

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Go GraphQL API with Echo!")
	})
	e.POST("/graphql", echo.WrapHandler(graphqlHandler))
	e.GET("/graphql", echo.WrapHandler(graphqlHandler))

	// Start the server
	e.Logger.Fatal(e.Start(":8999"))
}
