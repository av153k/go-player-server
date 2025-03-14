package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func DatabaseConnection() *pgx.Conn {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(".env file not found")

	}

	connUrl := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), connUrl)

	if err != nil {
		log.Fatal("Unable to open database connection")
	}

	return conn

}
