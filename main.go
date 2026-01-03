package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r, err := SetupServer()
	if err != nil {
		log.Fatal(err)
	}

	StartHTTPServer(r)
}
