package main

import (
	"log"
	"net/http"

	"example.com/poker"
)

const dbFileName = "game.db.json"

func main() {

	store, closeStore, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer closeStore()
	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen to server %v", err)
	}

}
