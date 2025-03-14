package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	db := DatabaseConnection()
	defer db.Close(context.Background())
	store := NewPostgresPlayerStore(db)
	server := PlayerServer{
		store: store,
	}
	player := "Abhishek"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertResponseStatusCode(t, response.Code, http.StatusOK)

	assertResponse(t, response.Body.String(), "3")
}