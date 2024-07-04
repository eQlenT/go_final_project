package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"

	_ "modernc.org/sqlite"
)

func Task(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch r.Method {
	// В нашем случае необходимо добавить обработку GET-запроса, который возвратит все параметры задачи по её идентификатору.
	// Если сейчас нажать на иконку редактирования задачи, появится ошибка.
	// Исправьте ситуацию — реализуйте обработчик GET-запроса /api/task?id=<идентификатор>.
	// Запрос должен возвращать JSON-объект со всеми полями задачи.
	case http.MethodGet:
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
		for rows.Next() {
			err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				utils.SendErr(w, err, http.StatusInternalServerError)
				return
			}
			if err = rows.Err(); err != nil {
				utils.SendErr(w, err, http.StatusInternalServerError)
				return
			}
		}
		response, err := json.Marshal(task)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(response)
	// Добавьте обработку PUT-запроса в хендлер для /api/task.
	// При этом нужно данные нужно проверять так же, как при добавлении задачи.
	// В случае успешного изменения должен возвращаться пустой JSON {}, а в случае ошибки, она записывается в поле error.
	case http.MethodPut:
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			utils.SendErr(w, err, http.StatusBadRequest)
			return
		}
		err = utils.CheckRequest(task)
		if err != nil {
			utils.SendErr(w, err, http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
			task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	case http.MethodPost:
		var id struct {
			ID int64 `json:"id"`
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

		nextDate, err := utils.CompleteRequest(request)
		if err != nil {
			utils.SendErr(w, err, http.StatusBadRequest)
			return
		}

		res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
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
	default:
		err := fmt.Errorf("no request method")
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
}