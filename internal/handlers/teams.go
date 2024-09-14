package handlers

import (
	"net/http"
	"weecal/internal/store/team"
	"weecal/web/templates"
)

func HandleListTeams(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := teamStore.ListTeams()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			Render(w, r, templates.TeamsError(), "Teams Error")
			return
		}

		Render(w, r, templates.ListTeams(teams), "Teams")
	}
}

func HandleCreateTeam(w http.ResponseWriter, r *http.Request) {
	err := templates.CreateTeam().Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
