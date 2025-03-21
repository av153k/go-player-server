package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) (winner string) {
	return strings.Replace(userInput, " wins", "", 1)
}


func NewCli(store PlayerStore, in io.Reader) *CLI {
	return &CLI{store, in}
}
