package handlers

import (
	"log/slog"
	"net/http"
	"weecal/internal/store/team"
	"weecal/web/templates"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func HandleListTeams(teamStore team.TeamStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := teamStore.ListTeams()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			templ.Handler(templates.TeamsError()).ServeHTTP(w, r)
			return
		}

		templ.Handler(templates.ListTeams(teams)).ServeHTTP(w, r)
	}
}

func HandleCreateTeamView() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		hxRequest := r.Header.Get("HX-Request")
		if hxRequest == "true" {
			templ.Handler(templates.CreateTeamComponent(team.TeamForm{}, map[string]string{})).ServeHTTP(w, r)
		} else {
			templ.Handler(templates.CreateTeam()).ServeHTTP(w, r)
		}

	}
}

func HandleViewTeam(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		teamId := chi.URLParam(r, "id")
		slog.Info("HandleViewTeam", "teamId", teamId)
		team, err := teamStore.ReadTeam(teamId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			templ.Handler(templates.TeamsError()).ServeHTTP(w, r)
			return
		}
		hxRequest := r.Header.Get("HX-Request")
		if hxRequest == "true" {
			templ.Handler(templates.ViewTeamComponent(team)).ServeHTTP(w, r)
		} else {
			templ.Handler(templates.ViewTeam(team)).ServeHTTP(w, r)
		}

	}
}

func HandleDeleteTeam(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		teamId := chi.URLParam(r, "id")
		slog.Info("HandleDeleteTeam", "teamId", teamId)
		err := teamStore.DeleteTeam(teamId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			templ.Handler(templates.TeamsError()).ServeHTTP(w, r)
			return
		}
		// w.Header().Set("HX-Redirect", "/teams")
		w.WriteHeader(http.StatusOK)

	}
}

func validateCreateTeam(teamForm team.TeamForm) map[string]string {
	validationErrors := map[string]string{}
	if teamForm.Name == "" {
		validationErrors["name"] = "Missing Name"
	}
	if teamForm.ShortName == "" {
		validationErrors["shortName"] = "Missing Short Name"
	}
	return validationErrors
}

func HandleCreateTeam(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		teamForm := team.TeamForm{
			Name:      r.FormValue("name"),
			ShortName: r.FormValue("shortName"),
		}
		validationErrors := validateCreateTeam(teamForm)
		if len(validationErrors) != 0 {
			templ.Handler(templates.CreateTeamComponent(teamForm, validationErrors)).ServeHTTP(w, r)
			return
		}

		team := team.Team{
			Name:      teamForm.Name,
			ShortName: teamForm.ShortName,
		}

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
