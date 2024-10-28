package team

import (
	"errors"
	"fmt"
	"log/slog"
	"weecal/internal/store/utils"

	"github.com/jmoiron/sqlx"
)

const TeamSchema = `
CREATE TABLE teams (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	short_name TEXT UNIQUE NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
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

	result, err := tx.NamedExec(`INSERT INTO teams (name, short_name, created_at, updated_at)
		VALUES (:name, :short_name, `+utils.DateTimeFormat+`, `+utils.DateTimeFormat+`);`, team)

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

	result, err := tx.NamedExec(`UPDATE teams SET name=:name, short_name=:short_name, updated_at=`+utils.DateTimeFormat+` WHERE id=:id AND updated_at=:updated_at;`, team)

	if err != nil {
		slog.Info("Error updating team", "err", err)
		_ = tx.Rollback()
		return err
	}
	if rows, _ := result.RowsAffected(); rows < 1 {
		_ = tx.Rollback()
		rereadTeam, err := s.ReadTeam(team.ID)
		if err != nil {
			slog.Info("Error updating team, not found", "id", rereadTeam.ID, "err", err)
		}
		if rereadTeam.UpdatedAt != team.UpdatedAt {
			slog.Info("Error updating team, already updated", "id", rereadTeam.ID, "err", err, "team.UpdatedAt", team.UpdatedAt, "rereadTeam.UpdatedAt", rereadTeam.UpdatedAt)
		}
		return errors.New(fmt.Sprintf("Error updating team with id: %d, not found or updated", team.ID))
	}

	err = tx.Commit()
	if err != nil {
		slog.Info("Error committing transaction", "err", err)
		return err
	}

	slog.Info("Updated team", "id", team.ID)

	return nil
}

// TODO: pass in team reference instead and return that
func (s *SQLTeamStore) ReadTeam(id int) (Team, error) {
	team := Team{
		ID: id,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		slog.Info("Error creating transaction", "err", err)
		return team, err
	}

	stmt, err := tx.PrepareNamed(`SELECT id, name, short_name, created_at, updated_at FROM teams WHERE id=:id;`)
	if err != nil {
		slog.Info("Error creating prepared statement to read team", "err", err, "id", id)
		_ = tx.Rollback()
		return team, err
	}
	defer stmt.Close()

	err = stmt.Get(&team, team)

	if err != nil {
		slog.Info("Error reading team", "err", err, "id", id)
		_ = tx.Rollback()
		return team, err
	}

	err = tx.Commit()
	if err != nil {
		slog.Info("Error committing transaction", "err", err)
		return team, err
	}

	return team, nil
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
