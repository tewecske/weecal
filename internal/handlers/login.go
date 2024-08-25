package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"weecal/internal/hash"
	"weecal/internal/store/session"
	"weecal/internal/store/user"
	"weecal/web/templates"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	Render(w, r, templates.Login(), "Login")
}

func HandlePostLogin(
	userStore user.UserStore,
	sessionStore session.SessionStore,
	passwordHash hash.PasswordHash,
	sessionCookieName string,
) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := userStore.GetUser(email)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			c := templates.LoginError()
			c.Render(r.Context(), w)
			return
		}
		slog.Info("Logging in user", "user email", email, "user", user.ID)

		passwordIsValid, err := passwordHash.ComparePasswordAndHash(password, user.Password)
		slog.Info("User password", "valid", passwordIsValid)

		if err != nil || !passwordIsValid {
			w.WriteHeader(http.StatusUnauthorized)
			c := templates.LoginError()
			c.Render(r.Context(), w)
			return
		}

		session, err := sessionStore.CreateSession(&session.Session{
			UserID: user.ID,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userID := user.ID
		sessionID := session.SessionID

		cookieValue := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", sessionID, userID)))

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     sessionCookieName,
			Value:    cookieValue,
			Expires:  expiration,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}
