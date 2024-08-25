package handlers

import (
	"fmt"
	"net/http"
	"time"
	"weecal/web/templates"
)

func HandleCalendar(w http.ResponseWriter, r *http.Request) {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	hours := []string{}
	for hour := 0; hour < 24; hour++ {
		hours = append(hours, fmt.Sprintf("%02d:00", hour))
	}

	Render(w, r, templates.Calendar(days, hours), "Calendar")
}

func Date(i1, i2, i3, i4, i5, i6, i7, i8 int, tz *time.Location) {
	panic("unimplemented")
}
