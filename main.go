package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	db := DatabaseConnection()
	defer db.Close(context.Background())
	store := NewPostgresPlayerStore(db)
	server := &PlayerServer{
		store: store,
	}
	log.Fatal(http.ListenAndServe(":5000", server))
}
