package httpserver

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("Unable to open db file %s: %v", dbFileName, err)
	}
	server := NewPlayerServer(&FileSystemPlayerStore{db})
	log.Fatal(http.ListenAndServe(":5000", server))
}
