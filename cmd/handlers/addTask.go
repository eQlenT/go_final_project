package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/cmd/utils"
	"go_final_project/models"
	"net/http"

	_ "modernc.org/sqlite"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	var errStr struct {
		Error string `json:"error"`
	}
	var id struct {
		ID int64 `json:"id"`
	}

	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		errStr.Error = "request method must be POST"
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
		return
	}

	var request models.Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errStr.Error = fmt.Sprintf("error decoding request body: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
		return
	}

	err = utils.CheckRequest(request)
	if err != nil {
		errStr.Error = fmt.Sprintf("error validating request: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		errStr.Error = fmt.Sprintf("error opening DB: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	nextDate, err := utils.CompleteRequest(request)
	if err != nil {
		errStr.Error = fmt.Sprintf("error processing request: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		nextDate, request.Title, request.Comment, request.Repeat)
	if err != nil {
		errStr.Error = fmt.Sprintf("error inserting into database: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		errStr.Error = fmt.Sprintf("error getting last insert ID: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}
	id.ID = lastInsertID
	response, err := json.Marshal(id)
	if err != nil {
		errStr.Error = fmt.Sprintf("error marshaling response: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(response)
}
