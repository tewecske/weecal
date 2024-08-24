package user

import (
	"weecal/internal/hash"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

const UserSchema = `
CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	email TEXT NOT NULL,
	password TEXT NOT NULL
);`

type SQLUserStore struct {
	db           *sqlx.DB
	passwordhash hash.PasswordHash
}

type NewUserStoreParams struct {
	DB           *sqlx.DB
	PasswordHash hash.PasswordHash
}

func NewUserStore(params NewUserStoreParams) *SQLUserStore {
	return &SQLUserStore{
		db:           params.DB,
		passwordhash: params.PasswordHash,
	}
}

func (s *SQLUserStore) CreateUser(email string, password string) error {

	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	result, err := s.db.NamedExec(`INSERT INTO users (email, password) VALUES (:email, :password);`,
		map[string]interface{}{
			"email":    email,
			"password": hashedPassword,
		})

	if err != nil {
		slog.Info("Error creating user", "err", err)
	} else {
		id, _ := result.LastInsertId()

		slog.Info("Inserting user", "lastId", id)
	}

	return err
}

func (s *SQLUserStore) GetUser(email string) (*User, error) {

	user := User{}
	rows, err := s.db.NamedQuery(`SELECT id, email, password FROM users where email=:email;`,
		map[string]interface{}{
			"email": email,
		})

	if rows.Next() {
		err = rows.StructScan(&user)
	}

	if err != nil {
		return nil, err
	}
	return &user, err
}
