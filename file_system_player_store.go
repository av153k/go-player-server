package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func initialisePlayerDbFile(file *os.File) error {
	file.Seek(0, io.SeekStart)
	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem loading player store from file %s, %w", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initialisePlayerDbFile(file)

	if err != nil {
		log.Fatalf("problem initializing player db file, %s", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %w", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		json.NewEncoder(&Tape{file}), league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (wins int) {
	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0

}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		f.league = append(league, Player{name, 1})
	}

	f.database.Encode(f.league)
}
