package handlers

import (
	"net/http"
	"time"
)

func HandlePostLogout(
	sessionCookieName string,
) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    sessionCookieName,
			MaxAge:  -1,
			Expires: time.Now().Add(-100 * time.Hour),
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
