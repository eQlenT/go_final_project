// Реализовать обработчик для GET-запроса /api/tasks.
// Он должен возвращать список ближайших задач в формате JSON в виде списка в поле tasks.
// Задачи должны быть отсортированы по дате в сторону увеличения.
// Каждая задача должна содержать все поля таблицы scheduler в виде строк.
// Дата представлена в уже знакомом вам формате 20060102.

package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/models"
	"net/http"

	_ "modernc.org/sqlite"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	var errStr models.MyErr
	var response struct {
		Tasks []models.Task `json:"tasks"`
	}
	tasks := make([]models.Task, 0, 50)

	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		errStr.Error = "request method must be GET"
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		errStr.Error = fmt.Sprint(err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM scheduler")
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		errStr.Error = fmt.Sprint(err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var (
			task    models.Task
			id      int
			title   string
			date    string
			repeat  string
			comment string
		)

		err := rows.Scan(&id, &date, &title, &comment, &repeat)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			errStr.Error = fmt.Sprint(err)
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
			return
		}
		task = models.Task{
			ID:      id,
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		tasks = append(tasks, models.Task{})
	}
	response.Tasks = tasks
	tasksJSON, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		errStr.Error = fmt.Sprint(err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(tasksJSON)
}
