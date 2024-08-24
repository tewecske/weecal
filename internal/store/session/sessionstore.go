package session

import (
	"fmt"
	"weecal/internal/store/user"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const SessionSchema = `
CREATE TABLE sessions (
	id INTEGER PRIMARY KEY,
	session_id TEXT UNIQUE,
	user_id INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id)
);`

type SQLSessionStore struct {
	db *sqlx.DB
}

type NewSessionStoreParams struct {
	DB *sqlx.DB
}

func NewSessionStore(params NewSessionStoreParams) *SQLSessionStore {
	return &SQLSessionStore{
		db: params.DB,
	}
}

func (s *SQLSessionStore) CreateSession(session *Session) (*Session, error) {

	session.SessionID = uuid.New().String()

	result, err := s.db.NamedExec(`INSERT INTO sessions (session_id, user_id) VALUES (:session_id, :user_id);`, map[string]interface{}{
		"session_id": session.SessionID,
		"user_id":    session.UserID,
	})

	if err != nil {
		slog.Info("Error creating session", "err", err)
		return nil, err
	}
	id, _ := result.LastInsertId()

	slog.Info("Inserting session", "lastId", id)

	return session, nil
}

func (s *SQLSessionStore) GetUserFromSession(sessionID string, userID string) (*user.User, error) {
	user := user.User{}
	rows, err := s.db.NamedQuery(`
		SELECT u.id as id, u.email as email, u.password as password
		FROM users as u
		JOIN sessions as s ON u.id = s.user_id
		WHERE s.session_id=:session_id AND u.id=:user_id`,
		map[string]interface{}{
			"session_id": sessionID,
			"user_id":    userID,
		})

	if rows.Next() {
		err = rows.StructScan(&user)
	}

	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("no user associated with the session")
	}

	return &user, nil
}
