package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"weecal/web/templates"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request)

func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if err := h(w, r); err != nil {
		// slog.Error("HTTP handler error", "error", err, "path", r.URL.Path)
		// }
		h(w, r)
	}
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	err := templates.Layout(c, "GOTH Starter").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
