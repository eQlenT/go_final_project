package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"

	_ "modernc.org/sqlite"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	var id struct {
		ID int64 `json:"id"`
	}

	if r.Method != http.MethodPost {
		err := fmt.Errorf("request method must be POST")
		utils.SendErr(w, err, http.StatusBadRequest)
		return
	}

	var request models.Task
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendErr(w, err, http.StatusBadRequest)
		return
	}

	err = utils.CheckRequest(request)
	if err != nil {
		utils.SendErr(w, err, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	nextDate, err := utils.CompleteRequest(request)
	if err != nil {
		utils.SendErr(w, err, http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		nextDate, request.Title, request.Comment, request.Repeat)
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	id.ID = lastInsertID
	response, err := json.Marshal(id)
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(response)
}
