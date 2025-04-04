package main

import (
	"fmt"
	"log"
	"os"

	"example.com/poker"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {name} wins to record a win")

	store, closeStore, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer closeStore()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.StdOutAlerter), store)

	poker.NewCLI(os.Stdin, os.Stdout, game).PlayPoker()

}
