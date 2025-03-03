package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {

	// Load the .env file
	_ = godotenv.Load(".env")

	// Read environment variables
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Connection successful")
	}
	defer conn.Close(context.Background())
}
