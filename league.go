package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}

	return nil
}

func NewLeague(rdr io.Reader) (league []Player, err error) {
	err = json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		return nil, fmt.Errorf("error while parsing league %v", err)
	}

	return league, nil
}
