package session

import "weecal/internal/store/user"

type TestSessionStore struct {
	Session *Session
	User    *user.User
	Error   error
}

func (s *TestSessionStore) CreateSession(session *Session) (*Session, error) {
	if s.Error != nil {
		return nil, s.Error
	} else {
		session.SessionID = s.Session.SessionID
		return session, nil
	}
}

func (s *TestSessionStore) GetUserFromSession(sessionID string, userID string) (*user.User, error) {
	if s.Error != nil {
		return nil, s.Error
	} else {
		return s.User, nil
	}
}
