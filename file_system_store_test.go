package httpserver

import (
	"io"
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

	store := FileSystemPlayerStore{database}

	t.Run("Get league", func(t *testing.T) {
		got := store.GetLeague()

		want := []Player{
			{"Chris", 10},
			{"Liam", 45},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()

		assertLeague(t, got, want)
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

}

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
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
