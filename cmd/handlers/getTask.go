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
	tasks := make(map[string][]models.Task)

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

	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
	ORDER BY date LIMIT 50`)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		errStr.Error = fmt.Sprint(err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			errStr.Error = fmt.Sprint(err)
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
			return
		}
		tasks["tasks"] = append(tasks["tasks"], task)
	}
	if tasks["tasks"] == nil {
		tasks["tasks"] = []models.Task{}
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		errStr.Error = fmt.Sprint(err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(response)
}
