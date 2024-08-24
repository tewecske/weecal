package handlers

import (
	"weecal/web/templates"
	"log/slog"
	"net/http"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	slog.Warn("Page not found", "path", r.URL.Path)
	Render(w, r, templates.NotFound())
}
