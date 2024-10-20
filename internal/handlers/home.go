package handlers

import (
	"net/http"
	"weecal/web/templates"

	"github.com/a-h/templ"
)

func HandleHome() http.HandlerFunc {
	return templ.Handler(templates.Index()).ServeHTTP
}
