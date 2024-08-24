package session

import "weecal/internal/store/user"

type Session struct {
	ID        uint
	SessionID string
	UserID    uint
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*user.User, error)
}
