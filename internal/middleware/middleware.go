package middleware

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"weecal/internal/store/session"
	"weecal/internal/store/user"
	u "weecal/internal/store/user"
	"net/http"
	"strings"
)

func TextHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

type AuthMiddleware struct {
	sessionStore      session.SessionStore
	sessionCookieName string
}

func NewAuthMiddleware(sessionStore session.SessionStore, sessionCookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore:      sessionStore,
		sessionCookieName: sessionCookieName,
	}
}

func (m *AuthMiddleware) AddUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionCookie, err := r.Cookie(m.sessionCookieName)

		if err != nil {
			fmt.Println("error getting session cookie", err)
			next.ServeHTTP(w, r)
			return
		}

		decodedValue, err := b64.StdEncoding.DecodeString(sessionCookie.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		splitValue := strings.Split(string(decodedValue), ":")

		if len(splitValue) != 2 {
			next.ServeHTTP(w, r)
			return
		}

		sessionID := splitValue[0]
		userID := splitValue[1]

		fmt.Println("sessionID", sessionID)
		fmt.Println("userID", userID)

		user, err := m.sessionStore.GetUserFromSession(sessionID, userID)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := u.NewContext(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *user.User {
	user, ok := user.FromContext(ctx)
	if user == nil || !ok {
		return nil
	}

	return user
}
