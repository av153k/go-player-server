package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	db, closeDb := createTempFile(t, "[]")
	defer closeDb()
	store, err := NewFileSystemPlayerStore(db)

	assertNoError(t, err)
	server := NewPlayerServer(store)
	player := "Abhishek"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertResponseStatusCode(t, response.Code, http.StatusOK)

		assertResponse(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertContentTypeHeader(t, *response, "application/json")

		got := getLeagueFromResponse(t, response.Body)

		want := []Player{
			{player, 3},
		}

		assertLeague(t, got, want)
	})
}
