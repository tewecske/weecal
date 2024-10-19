package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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

func HandleCreateTeamView() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hxRequest := r.Header.Get("HX-Request")
		if hxRequest == "true" {
			err := templates.CreateTeamComponent().Render(r.Context(), w)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
		} else {
			Render(w, r, templates.CreateTeam(), "Teams")
		}

	}
}

func HandleCreateTeam(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var team team.Team
		r.ParseForm()
		name := r.FormValue("name")
		shortName := r.FormValue("shortName")
		var validationErrors []string = make([]string, 0, 2)
		if name == "" {
			validationErrors = append(validationErrors, "Missing Name")
		}
		if shortName == "" {
			validationErrors = append(validationErrors, "Missing Short Name")
		}
		if len(validationErrors) != 0 {
			errors := strings.Join(validationErrors, ",")
			http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusInternalServerError)
			return
		}

		team.Name = name
		team.ShortName = shortName

		slog.Info("Decoded team from request", "team", team)
		err := teamStore.CreateTeam(team)
		if err != nil {
			http.Error(w, "Error creating team", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/teams")
		w.WriteHeader(http.StatusOK)
	}
}
