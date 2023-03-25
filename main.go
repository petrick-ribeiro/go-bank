package main

import (
	"log"

	"github.com/petrick-ribeiro/go-bank/api"
	"github.com/petrick-ribeiro/go-bank/storage"
)

// TODO: Dockerize all stuff
func main() {
	store, err := storage.NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// TODO: Create a flag to server port
	server := api.NewAPIServer(":3000", store)
	server.Run()
}
