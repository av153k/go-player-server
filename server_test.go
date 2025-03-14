package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"abhishek": 20,
			"damon":    10,
		},
	}
	server := &PlayerServer{
		store: &store,
	}

	t.Run("Return abhishek's score", func(t *testing.T) {
		request := newGetScoreRequest("abhishek")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertResponse(t, got, want)

	})

	t.Run("Return Damon's score", func(t *testing.T) {
		request := newGetScoreRequest("damon")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertResponse(t, got, want)
	})

	t.Run("return 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("julia")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertResponseStatusCode(t, got, want)
	})
}

func newGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func assertResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Got: %q, Want: %q", got, want)
	}
}

func assertResponseStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("StatusCode: Got %d want %d", got, want)
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,

	}
	server := &PlayerServer{
		store: &store,
	}

	t.Run("it returns accepted on post", func(t *testing.T) {
		player := "abhishek"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}


		if store.winCalls[0] != player {
			t.Errorf("Did not get correct player, Got %s want %s", store.winCalls[0], player)
		}

	})
}
