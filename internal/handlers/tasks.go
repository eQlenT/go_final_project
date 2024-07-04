// Реализовать обработчик для GET-запроса /api/tasks.
// Он должен возвращать список ближайших задач в формате JSON в виде списка в поле tasks.
// Задачи должны быть отсортированы по дате в сторону увеличения.
// Каждая задача должна содержать все поля таблицы scheduler в виде строк.
// Дата представлена в уже знакомом вам формате 20060102.

package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := make(map[string][]models.Task)
	search := r.FormValue("search")
	isSearch := search != ""

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var rows *sql.Rows
	if isSearch {
		_, err := time.Parse("02.01.2006", search)
		if err != nil {
			rows, err = db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
	WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT 25`,
				sql.Named("search", "%"+search+"%"))
			if err != nil {
				utils.SendErr(w, err, http.StatusInternalServerError)
			}
		} else {
			date, _ := time.Parse("02.01.2006", search)
			dateFormat := date.Format("20060102")
			rows, err = db.Query(`SELECT id, date, title, comment, repeat FROM scheduler
		WHERE date = :date LIMIT 25`,
				sql.Named("date", dateFormat))
			if err != nil {
				utils.SendErr(w, err, http.StatusInternalServerError)
			}
		}
	} else {
		rows, err = db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
	ORDER BY date LIMIT 25`)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}

		tasks["tasks"] = append(tasks["tasks"], task)
	}
	if err := rows.Err(); err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	// Если задач нет, возвращаем пустой json.
	if tasks["tasks"] == nil {
		tasks["tasks"] = []models.Task{}
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(response)
}
