package handlers

import (
	"weecal/web/templates"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	Render(w, r, templates.Index())
}
