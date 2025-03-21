package poker

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	database, cleanDatabase := createTempFile(t, `
		[
			{"name": "Chris", "wins": 10},
			{"name": "Liam", "wins": 45}
		]`)

	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)

	assertNoError(t, err)

	t.Run("Get league", func(t *testing.T) {
		got := store.GetLeague()

		want := []Player{
			{"Liam", 45},
			{"Chris", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()

		assertLeague(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDb := createTempFile(t, "")

		defer cleanDb()

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

	})

	t.Run("Get player score", func(t *testing.T) {
		got := store.GetPlayerScore("Chris")

		want := 10

		assertPlayerScore(t, got, want)
	})

	t.Run("Record wins", func(t *testing.T) {
		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")

		assertPlayerScore(t, got, 11)
	})

	t.Run("Record wins for new player", func(t *testing.T) {
		store.RecordWin("Abhishek")
		got := store.GetPlayerScore("Abhishek")

		assertPlayerScore(t, got, 1)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `
		[{"Name":"Subhadeep", "Wins":10},
		{"Name":"Samrat", "Wins":30}
		]`)

		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{"Samrat", 30},
			{"Subhadeep", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()

		assertLeague(t, got, want)
	})

}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("Unable to create temp file %v", err)
	}

	tempFile.Write([]byte(initialData))

	removFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removFile
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)

	}
}
