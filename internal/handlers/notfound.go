package handlers

import (
	"net/http"
	"weecal/web/templates"

	"github.com/a-h/templ"
)

func HandleNotFound() http.HandlerFunc {
	// slog.Warn("Page not found", "path", r.URL.Path)
	return templ.Handler(templates.NotFound()).ServeHTTP
}
