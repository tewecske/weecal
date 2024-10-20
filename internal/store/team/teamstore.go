package team

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

const TeamSchema = `
CREATE TABLE teams (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	short_name TEXT UNIQUE NOT NULL
);`

type SQLTeamStore struct {
	db *sqlx.DB
}

type NewTeamStoreParams struct {
	DB *sqlx.DB
}

func NewTeamStore(params NewTeamStoreParams) *SQLTeamStore {
	return &SQLTeamStore{
		db: params.DB,
	}
}

func (s *SQLTeamStore) ListTeams() ([]Team, error) {
	teams := []Team{}
	err := s.db.Select(&teams, `SELECT id, name, short_name FROM teams;`)

	if err != nil {
		slog.Info("Error listing teams", "err", err)
		return nil, err
	}
	slog.Info("Got teams from database", "teams", teams)

	return teams, nil
}

func (s *SQLTeamStore) CreateTeam(team Team) error {

	result, err := s.db.NamedExec(`INSERT INTO teams (name, short_name) VALUES (:name, :short_name);`, map[string]interface{}{
		"name":       team.Name,
		"short_name": team.ShortName,
	})

	if err != nil {
		slog.Info("Error creating team", "err", err)
		return err
	}
	id, _ := result.LastInsertId()

	slog.Info("Inserting team", "lastId", id)

	return nil
}

func (s *SQLTeamStore) ReadTeam(id string) (Team, error) {
	team := Team{}
	rows, err := s.db.NamedQuery(`SELECT id, name, short_name FROM teams WHERE id=:id;`, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		slog.Info("Error reading teams", "err", err, "id", id)
		return team, err
	}
	if rows.Next() {
		rows.StructScan(&team)
		slog.Info("Got team from database", "team", team)

		return team, nil
	}
	return team, errors.New(fmt.Sprintf("Error reading team with id: %s", id))
}
