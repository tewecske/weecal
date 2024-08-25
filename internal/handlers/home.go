package handlers

import (
	"net/http"
	"weecal/web/templates"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	Render(w, r, templates.Index(), "Home")
}
