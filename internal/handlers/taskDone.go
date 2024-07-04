package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"
	"time"
)

func TaskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := errors.New("request method must be post")
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	id := r.FormValue("id")
	if id == "" {
		utils.SendErr(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}
	task := models.Task{}
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
	}
	if err = rows.Err(); err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	if task.Repeat != "" {
		task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`,
			task.Date, task.ID)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	} else {
		_, err = db.Exec(`DELETE FROM scheduler WHERE id = ?`, task.ID)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	}
}
