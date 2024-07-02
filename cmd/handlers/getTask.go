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
	tasks := make(map[string][]models.Task)

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println("open db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler 
	ORDER BY date`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tasks["tasks"] = append(tasks["tasks"], task)

	}

	// Если задач нет, возвращаем пустой json.
	if tasks["tasks"] == nil {
		tasks["tasks"] = []models.Task{}
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
