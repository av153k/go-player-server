package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)


type CLI struct {
	playerStore  PlayerStore
	in           *bufio.Scanner
	blindAlerter BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, blindAlerter BlindAlerter) *CLI {
	return &CLI{store, bufio.NewScanner(
		in,
	), blindAlerter}
}

func (cli *CLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.blindAlerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Second
	}
}

func extractWinner(userInput string) (winner string) {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
