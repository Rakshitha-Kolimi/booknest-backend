package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r, err := setupServer()
	if err != nil {
		log.Fatal(err)
	}

	startHTTPServer(r)
}
