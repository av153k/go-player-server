package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})
		request := newGameRequest()

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, 200)

	})

	t.Run("when we get a message over a websocket it is the winner of the game", func(t *testing.T) {
		store := &StubPlayerStore{}
		winner := "Ruth"
		server := httptest.NewServer(NewPlayerServer(store))

		defer server.Close()

		wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)

		if err != nil {
			t.Fatalf("Unable to open a websocket connection: %s %v", wsUrl, err)
		}

		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
			t.Fatalf("could not send message over the websocket connection %v", err)
		}


		AssertPlayerWin(t, store, winner)

	})
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"abhishek": 20,
			"damon":    10,
		},
	}
	server := NewPlayerServer(&store)

	t.Run("Return abhishek's score", func(t *testing.T) {
		request := newGetScoreRequest("abhishek")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertStatus(t, response, http.StatusOK)
		assertResponse(t, got, want)

	})

	t.Run("Return Damon's score", func(t *testing.T) {
		request := newGetScoreRequest("damon")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertStatus(t, response, http.StatusOK)
		assertResponse(t, got, want)
	})

	t.Run("return 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("julia")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := http.StatusNotFound

		assertStatus(t, response, want)
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

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func assertResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Got: %q, Want: %q", got, want)
	}
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Errorf("StatusCode: Got %d want %d", response.Code, want)
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("it returns accepted on post", func(t *testing.T) {
		player := "abhishek"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertPlayerWin(t, &store, "abhishek")

	})
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse json from the response body %q into slice of player, %v", response.Body, err)
		}

		assertStatus(t, response, http.StatusOK)

	})

	t.Run("it returns the league table as json", func(t *testing.T) {
		wantedPlayers := []Player{
			{"chris", 3},
			{"liam", 6},
			{"leo", 2},
		}

		store := StubPlayerStore{nil, nil, wantedPlayers}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertContentTypeHeader(t, *response, "application/json")
		assertStatus(t, response, http.StatusOK)
		assertLeague(t, got, wantedPlayers)
	})
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse json server body %v to slice of players: %v", body, err)

	}

	return
}

func assertContentTypeHeader(t testing.TB, response httptest.ResponseRecorder, contentType string) {
	t.Helper()
	if response.Result().Header.Get("Content-Type") != contentType {
		t.Errorf("Header did not have content-type as %q: got %v", contentType, response.Result().Header)
	}
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Players mismatch: got %v want %v", got, want)
	}
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func assertPlayerScore(t testing.TB, got, want int) {
	if got != want {
		t.Errorf("player score mismatched: got %q want %q", got, want)
	}
}
