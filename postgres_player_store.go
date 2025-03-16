package httpserver

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewPostgresPlayerStore(db *pgx.Conn) *PostgresPlayerStore {
	return &PostgresPlayerStore{db}
}

type PostgresPlayerStore struct {
	db *pgx.Conn
}

func (p *PostgresPlayerStore) GetPlayerScore(name string) int {
	var id int
	var playerName string
	var score int
	err := p.db.QueryRow(context.Background(), `SELECT id, name, score from "go-player-server".player where name=$1`, name).Scan(&id, &playerName, &score)

	if err != nil {
		fmt.Println(fmt.Errorf("player not found: %w", err))
		return 0
	}

	return score
}

func (p *PostgresPlayerStore) RecordWin(name string) {
	score := p.GetPlayerScore(name)
	var userId int
	if score == 0 {
		query := `INSERT INTO "go-player-server".player (name, score) VALUES ($1, $2) RETURNING id`
		err := p.db.QueryRow(context.Background(), query, name, 1).Scan(&userId)
		if err != nil {
			fmt.Println(fmt.Errorf("unable to record win: %w", err))
		}
	} else {
		query := `UPDATE "go-player-server".player SET score=$1 where name=$2`
		_, err := p.db.Exec(context.Background(), query, score+1, name)
		if err != nil {
			fmt.Println(fmt.Errorf("unable to record win: %w", err))
		}
	}
}

func (p *PostgresPlayerStore) GetLeague() League {
	return []Player{}
}
