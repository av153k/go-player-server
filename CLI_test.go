package poker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"example.com/poker"
)

type scheduledAlert struct {
	scheduledAt time.Duration
	amount      int
}

func (s *scheduledAlert) string() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.scheduledAt)
}

type SpyBindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{
		duration, amount,
	})
}

func TestCLI(t *testing.T) {

	var dummyBlindAlerter = &SpyBindAlerter{}

	t.Run("record Chris win from user input", func(t *testing.T) {

		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in, dummyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record Samrat win from user input", func(t *testing.T) {
		in := strings.NewReader("Samrat wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in, dummyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Samrat")
	})

	t.Run("it schedules printing of bind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(want.string(), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]

				assertScheduledAlert(t, got, want)
			})
		}

	})

}

func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	t.Helper()

	amountGot := got.amount

	if amountGot != want.amount {
		t.Errorf("got amount %d, want %d", amountGot, want.amount)
	}

	gotScheduledTime := got.scheduledAt

	if gotScheduledTime != want.scheduledAt {
		t.Errorf("got scheduled time %v, want %v", gotScheduledTime, want.scheduledAt)
	}
}
