package poker_test

import (
	"fmt"
	"testing"
	"time"

	"example.com/poker"
)

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(5)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Second, Amount: 200},
			{At: 20 * time.Second, Amount: 300},
			{At: 30 * time.Second, Amount: 400},
			{At: 40 * time.Second, Amount: 500},
			{At: 50 * time.Second, Amount: 600},
			{At: 60 * time.Second, Amount: 800},
			{At: 70 * time.Second, Amount: 1000},
			{At: 80 * time.Second, Amount: 2000},
			{At: 90 * time.Second, Amount: 4000},
			{At: 100 * time.Second, Amount: 8000},
		}

		checkSchedulingCases(t, cases, blindAlerter)

	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(7)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Second, Amount: 200},
			{At: 24 * time.Second, Amount: 300},
			{At: 36 * time.Second, Amount: 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewTexasHoldem(dummyBlindAlerter, store)
	winner := "Subha"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}


func checkSchedulingCases(t *testing.T, cases []poker.ScheduledAlert, blindAlerter *poker.SpyBindAlerter) {
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(blindAlerter.Alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.Alerts)
			}

			got := blindAlerter.Alerts[i]

			poker.AssertScheduledAlert(t, got, want)
		})
	}
}