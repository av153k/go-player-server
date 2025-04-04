package poker

import (
	"fmt"
	"testing"
	"time"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

type ScheduledAlert struct {
	At time.Duration
	Amount      int
}

func (s *ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyBindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.Alerts = append(s.Alerts, ScheduledAlert{
		duration, amount,
	})
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {

	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want  %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Fatalf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}



func AssertScheduledAlert(t testing.TB, got, want ScheduledAlert) {
	t.Helper()

	amountGot := got.Amount

	if amountGot != want.Amount {
		t.Errorf("got amount %d, want %d", amountGot, want.Amount)
	}

	gotScheduledTime := got.At

	if gotScheduledTime != want.At {
		t.Errorf("got scheduled time %v, want %v", gotScheduledTime, want.At)
	}
}
