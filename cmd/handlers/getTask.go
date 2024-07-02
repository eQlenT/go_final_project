// Реализовать обработчик для GET-запроса /api/tasks.
// Он должен возвращать список ближайших задач в формате JSON в виде списка в поле tasks.
// Задачи должны быть отсортированы по дате в сторону увеличения.
// Каждая задача должна содержать все поля таблицы scheduler в виде строк.
// Дата представлена в уже знакомом вам формате 20060102.

package handlers

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/models"
	"net/http"
	"slices"
	"strconv"

	_ "modernc.org/sqlite"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	var errStr models.MyErr
	type taskMap map[string]string
	var response struct {
		Tasks []taskMap `json:"tasks"`
	}
	tasks := make([]taskMap, 0, 50)

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
			task    taskMap
			id      int
			title   string
			date    string
			repeat  string
			comment string
		)
		/*
			m := make(map[string]string)
			m["id"] = "2"
			m["title"] = "complete project"
			m["date"] = "20240703"
			m["repeat"] = ""
			m["comment"] = ""
			fmt.Println(m)

			// convert map to json
			jsonString, _ := json.Marshal(m)
			fmt.Println(string(jsonString))

			// convert json to struct
			s := Task{}
			json.Unmarshal(jsonString, &s)
			fmt.Println(s)
		*/

		err := rows.Scan(&id, &date, &title, &comment, &repeat)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			errStr.Error = fmt.Sprint(err)
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errStr.Error), http.StatusBadRequest)
			return
		}
		idString := strconv.Itoa(id)
		task = taskMap{
			"id":      idString,
			"date":    date,
			"title":   title,
			"comment": comment,
			"repeat":  repeat,
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		tasks = append(tasks, taskMap{"": ""})
	}

	slices.SortFunc(tasks, func(a, b taskMap) int {
		if n := cmp.Compare(a["date"], b["date"]); n != 0 {
			return n
		}
		return 0
	})
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
