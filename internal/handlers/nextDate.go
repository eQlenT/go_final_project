package handlers

import (
	"fmt"
	"go_final_project/internal/utils"
	"net/http"
	"time"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\nневерный формат now", err), http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if date == "" || repeat == "" {
		http.Error(w, fmt.Sprintf("%s\nневерный формат date или repeat", err), http.StatusBadRequest)
		return
	}
	next, err := utils.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", next)
}
