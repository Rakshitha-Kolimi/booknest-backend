package main

import (
	"log"

	"github.com/joho/godotenv"

	"booknest/internal/http/database"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot load the env file", "err", err)
	}

	// connect to database
	dbpool, err := database.Connect()
	if err != nil {
		log.Fatal("Cannot connect to database", "err", err)
	}

	defer dbpool.Close()

	// Set up server
	r, err := SetupServer(dbpool)
	if err != nil {
		log.Fatal(err)
	}

	StartHTTPServer(r)
}
