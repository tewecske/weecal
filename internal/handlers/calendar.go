package handlers

import (
	"fmt"
	"net/http"
	"weecal/web/templates"

	"github.com/a-h/templ"
)

func HandleCalendar() http.HandlerFunc {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	hours := []string{}
	for hour := 0; hour < 24; hour++ {
		hours = append(hours, fmt.Sprintf("%02d:00", hour))
	}

	return templ.Handler(templates.Calendar(days, hours)).ServeHTTP
}
