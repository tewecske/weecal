package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"weecal/internal/store/team"
	"weecal/web/templates"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/mattn/go-sqlite3"
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
			templ.Handler(templates.CreateTeam(team.TeamForm{}, map[string]string{})).ServeHTTP(w, r)
		}

	}
}

func HandleUpdateTeamView(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		paramId := chi.URLParam(r, "id")
		teamId, err := strconv.Atoi(paramId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templ.Handler(templates.TeamsError()).ServeHTTP(w, r)
			return
		}
		slog.Info("HandleUpdateTeam", "teamId", teamId)
		teamData, err := teamStore.ReadTeam(teamId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			templ.Handler(templates.TeamsError()).ServeHTTP(w, r)
			return
		}
		hxRequest := r.Header.Get("HX-Request")
		if hxRequest == "true" {
			templ.Handler(templates.UpdateTeamComponent(team.NewTeamForm(teamData), map[string]string{})).ServeHTTP(w, r)
		} else {
			templ.Handler(templates.UpdateTeam(team.NewTeamForm(teamData), map[string]string{})).ServeHTTP(w, r)
		}
	}
}

func validateUpdateTeam(teamForm team.TeamForm) map[string]string {
	validationErrors := map[string]string{}
	if teamForm.Name == "" {
		validationErrors["name"] = "Missing Name"
	}
	if teamForm.ShortName == "" {
		validationErrors["shortName"] = "Missing Short Name"
	}
	return validationErrors
}

func HandleUpdateTeam(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		teamForm := team.TeamForm{
			Name:      r.FormValue("name"),
			ShortName: r.FormValue("shortName"),
			UpdatedAt: r.FormValue("updatedAt"),
		}
		pathId := r.PathValue("id")
		id, err := strconv.Atoi(pathId)
		if err != nil {
			templ.Handler(templates.UpdateTeamComponent(teamForm,
				map[string]string{"globalError": fmt.Sprintf("Team not found with id: %s", pathId)})).ServeHTTP(w, r)
			return
		}
		teamForm.ID = id
		validationErrors := validateUpdateTeam(teamForm)
		if len(validationErrors) != 0 {
			w.Header().Set("HX-Reswap", "outerHTML")
			w.WriteHeader(http.StatusBadRequest)
			templ.Handler(templates.UpdateTeamComponent(teamForm, validationErrors)).ServeHTTP(w, r)
			return
		}

		team := team.Team{
			ID:        teamForm.ID,
			Name:      teamForm.Name,
			ShortName: teamForm.ShortName,
			UpdatedAt: teamForm.UpdatedAt,
		}

		slog.Info("Decoded team from request", "team", team)
		err = teamStore.UpdateTeam(team)
		if err != nil {
			w.Header().Set("HX-Reswap", "outerHTML")
			w.WriteHeader(http.StatusConflict)
			templ.Handler(templates.UpdateTeamComponent(teamForm,
				map[string]string{"globalError": fmt.Sprintf("Team with id: %s was already updated!", pathId)})).ServeHTTP(w, r)
			return
		}

		w.Header().Set("HX-Redirect", "/teams")
		w.WriteHeader(http.StatusOK)
	}
}
func HandleViewTeam(teamStore team.TeamStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		paramId := chi.URLParam(r, "id")
		teamId, err := strconv.Atoi(paramId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templ.Handler(templates.TeamsError()).ServeHTTP(w, r)
			return
		}
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
		paramId := chi.URLParam(r, "id")
		teamId, err := strconv.Atoi(paramId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templ.Handler(templates.TeamsError()).ServeHTTP(w, r)
			return
		}
		slog.Info("HandleDeleteTeam", "teamId", teamId)
		err = teamStore.DeleteTeam(teamId)
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

func fieldFromColumn(s any, f string) (string, error) {
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag
		db := tag.Get("db")
		if db == f {
			return field.Name, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Invalid field: %s", f))
}
func teamError(w http.ResponseWriter, err error) (map[string]string, error) {
	if sqerr, ok := err.(sqlite3.Error); ok {
		slog.Info("sqlite3 error", "sqerr", sqerr)
		// sqlite3 error: sqlite3.Error{Code:19, ExtendedCode:2067, SystemErrno:0x0, err:"UNIQUE constraint failed: teams.short_name"}
		switch sqerr.Code {
		case 19:
			switch sqerr.ExtendedCode {
			case 2067:
				words := strings.Split(sqerr.Error(), ".")
				columnName := words[len(words)-1]
				slog.Info("column", "columnName", columnName)
				fieldName, ferr := fieldFromColumn(team.Team{}, columnName)
				fieldName = strings.ToLower(fieldName[:1]) + fieldName[1:]
				slog.Info("field", "fieldName", fieldName)
				if ferr != nil {
					http.Error(w, "Error!", http.StatusInternalServerError)
				}
				// TODO: retarget for proper field
				validationErrors := map[string]string{}
				validationErrors[fieldName] = fmt.Sprintf("Duplicate %s!", fieldName)
				return validationErrors, nil
			default:
				http.Error(w, "Error!", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Error!", http.StatusInternalServerError)
		}
	} else {
		slog.Error("Unknown error!", "err", err)
		http.Error(w, "Error!", http.StatusInternalServerError)
	}
	return nil, errors.New("Error!")
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
			w.Header().Set("HX-Reswap", "outerHTML")
			w.WriteHeader(http.StatusBadRequest)
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
			validationErrors, _ = teamError(w, err)
			slog.Info("Validation errors", "validationErrors", validationErrors)
			w.Header().Set("HX-Reswap", "outerHTML")
			w.WriteHeader(http.StatusBadRequest)
			templ.Handler(templates.CreateTeamComponent(teamForm, validationErrors)).ServeHTTP(w, r)
			return
		}

		w.Header().Set("HX-Redirect", "/teams")
		w.WriteHeader(http.StatusOK)
	}
}
