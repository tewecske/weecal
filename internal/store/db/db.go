package db

import (
	"log/slog"
	"os"
	"weecal/internal/hash"
	"weecal/internal/store/session"
	"weecal/internal/store/team"
	"weecal/internal/store/user"

	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"

	"github.com/qustavo/sqlhooks/v2"
)

type DBAccess struct {
	DB           *sqlx.DB
	UserStore    user.UserStore
	SessionStore session.SessionStore
	TeamStore    team.TeamStore
}

func SetupDB(dbName string, passwordHash hash.PasswordHash) *DBAccess {
	db := Connect(dbName)

	userStore := user.NewUserStore(user.NewUserStoreParams{
		DB:           db,
		PasswordHash: passwordHash,
	})

	sessionStore := session.NewSessionStore(session.NewSessionStoreParams{
		DB: db,
	})

	teamStore := team.NewTeamStore(team.NewTeamStoreParams{
		DB: db,
	})

	return &DBAccess{
		DB:           db,
		UserStore:    userStore,
		SessionStore: sessionStore,
		TeamStore:    teamStore,
	}
}

// Hooks satisfies the sqlhook.Hooks interface
type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	fmt.Printf("> %s %q", query, args)
	return context.WithValue(ctx, "begin", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)
	fmt.Printf(". took: %s\n", time.Since(begin))
	return ctx, nil
}

func Connect(dbName string) *sqlx.DB {

	MustCreateTmp()

	// First, register the wrapper
	sql.Register("sqlite3WithHooks", sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, &Hooks{}))

	// Connect to the registered wrapped driver
	// db, _ := sql.Open("sqlite3WithHooks", ":memory:")
	db := sqlx.MustConnect("sqlite3WithHooks", dbName)

	// TODO Move to parameter
	userResult, err := db.MustExec(user.UserSchema).LastInsertId()
	if err != nil {
		panic(err)
	}
	slog.Info("User schema creation result", "userResult", userResult)
	sessionResult, err := db.MustExec(session.SessionSchema).LastInsertId()
	if err != nil {
		panic(err)
	}
	slog.Info("Session schema creation result", "sessionResult", sessionResult)

	teamResult, err := db.MustExec(team.TeamSchema).LastInsertId()
	if err != nil {
		panic(err)
	}
	slog.Info("Team schema creation result", "teamResult", teamResult)

	return db
}

func MustCreateTmp() {
	err := os.MkdirAll("/tmp", 0755)
	if err != nil {
		panic(err)
	}
}
