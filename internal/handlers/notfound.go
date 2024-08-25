package handlers

import (
	"log/slog"
	"net/http"
	"weecal/web/templates"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	slog.Warn("Page not found", "path", r.URL.Path)
	Render(w, r, templates.NotFound(), "Not Found!")
}
