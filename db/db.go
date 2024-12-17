package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"os"
)

var Pool *pgxpool.Pool

func Init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	Pool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}

	fmt.Println("Connected to PostgreSQL database")
}
