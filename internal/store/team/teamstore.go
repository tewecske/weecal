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
	tx, err := s.db.Beginx()
	if err != nil {
		slog.Info("Error creating transaction", "err", err)
		return teams, err
	}

	err = tx.Select(&teams, `SELECT id, name, short_name FROM teams;`)

	if err != nil {
		slog.Info("Error listing teams", "err", err)
		_ = tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		slog.Info("Error committing transaction", "err", err)
		return teams, err
	}

	slog.Info("Got teams from database", "teams", teams)

	return teams, nil
}

func (s *SQLTeamStore) CreateTeam(team Team) error {

	tx, err := s.db.Beginx()
	if err != nil {
		slog.Info("Error creating transaction", "err", err)
		return err
	}

	result, err := tx.NamedExec(`INSERT INTO teams (name, short_name) VALUES (:name, :short_name);`, team)

	if err != nil {
		slog.Info("Error creating team", "err", err)
		_ = tx.Rollback()
		return err
	}
	id, _ := result.LastInsertId()

	err = tx.Commit()
	if err != nil {
		slog.Info("Error committing transaction", "err", err)
		return err
	}
	slog.Info("Inserted team", "lastId", id)

	return nil
}

func (s *SQLTeamStore) UpdateTeam(team Team) error {

	tx, err := s.db.Beginx()
	if err != nil {
		slog.Info("Error creating transaction", "err", err)
		return err
	}

	result, err := tx.NamedExec(`UPDATE teams SET name=:name, short_name=:short_name WHERE id=:id;`, team)

	if err != nil {
		slog.Info("Error updating team", "err", err)
		_ = tx.Rollback()
		return err
	}
	if rows, _ := result.RowsAffected(); rows < 1 {
		slog.Info("Error updating team, not found", "id", team.ID)
		_ = tx.Rollback()
		return errors.New(fmt.Sprintf("Error updating team with id: %d, not found", team.ID))
	}

	err = tx.Commit()
	if err != nil {
		slog.Info("Error committing transaction", "err", err)
		return err
	}

	slog.Info("Updated team", "id", team.ID)

	return nil
}

func (s *SQLTeamStore) ReadTeam(id int) (Team, error) {
	team := Team{
		ID: id,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		slog.Info("Error creating transaction", "err", err)
		return team, err
	}

	stmt, err := tx.PrepareNamed(`SELECT id, name, short_name FROM teams WHERE id=:id;`)
	if err != nil {
		slog.Info("Error creating prepared statement to read team", "err", err, "id", id)
		return team, err
	}
	defer stmt.Close()

	err = stmt.Get(&team, team)

	if err != nil {
		slog.Info("Error reading team", "err", err, "id", id)
		return team, err
	}
	return team, errors.New(fmt.Sprintf("Error reading team with id: %d", id))
}

func (s *SQLTeamStore) DeleteTeam(id int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		slog.Info("Error creating transaction", "err", err)
		return err
	}

	result, err := tx.NamedExec(`DELETE FROM teams WHERE id=:id;`, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		slog.Info("Error deleting team", "err", err, "id", id)
		_ = tx.Rollback()
		return err
	}

	if rows, _ := result.RowsAffected(); rows < 1 {
		slog.Info("Error deleting team, not found", "id", id)
		_ = tx.Rollback()
		return errors.New(fmt.Sprintf("Error deleting team with id: %d, not found", id))
	}

	err = tx.Commit()
	if err != nil {
		slog.Info("Error committing transaction", "err", err)
		return err
	}

	return nil
}
