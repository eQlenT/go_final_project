package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

// Task обрабатывает HTTP-запросы для выполнения CRUD-операций над задачами в приложении-планировщике.
// Он поддерживает методы GET, POST, PUT и DELETE.
//
// GET:
//   - Извлекает задачу по её ID из базы данных.
//   - Возвращает JSON-объект со всеми полями задачи.
//   - Если ID пуст или не найден, возвращает ошибку 400 Bad Request.
//
// POST:
//   - Создает новую задачу в базе данных.
//   - Возвращает JSON-объект с ID созданной задачи.
//   - Если тело запроса недействительно или отсутствуют поля, возвращает ошибку 400 Bad Request.
//
// PUT:
//   - Обновляет существующую задачу в базе данных.
//   - Возвращает пустой JSON-объект.
//   - Если тело запроса недействительно или отсутствуют поля, возвращает ошибку 400 Bad Request.
//   - Если указанный ID не найден в базе данных, возвращает ошибку 400 Bad Request.
//
// DELETE:
//   - Удаляет задачу из базы данных по её ID.
//   - Возвращает пустой JSON-объект.
//   - Если ID пуст или не найден, возвращает ошибку 400 Bad Request.
//
// По умолчанию, если метод запроса не GET, POST, PUT или DELETE, возвращает ошибку 500 Internal Server Error.
func (c *DBConnection) Task(w http.ResponseWriter, r *http.Request) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
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
		rows, err := c.DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
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
		if task.ID == "" && task.Date == "" && task.Title == "" && task.Repeat == "" && task.Comment == "" {
			err = fmt.Errorf("no rows for id %s", id)
			utils.SendErr(w, err, http.StatusBadRequest)
			return
		}
		response, err := json.Marshal(task)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(response)

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
		} else {
			var maxID int
			row := c.DB.QueryRow(`SELECT MAX(id) FROM scheduler`)
			row.Scan(&maxID)
			if err = row.Err(); err != nil {
				utils.SendErr(w, err, http.StatusInternalServerError)
				return
			}
			givenID, err := strconv.Atoi(task.ID)
			if err != nil {
				err = errors.New("can not parse ID")
				utils.SendErr(w, err, http.StatusInternalServerError)
				return
			}
			if givenID > maxID {
				err = errors.New("given ID is more than number of rows in DB")
				utils.SendErr(w, err, http.StatusBadRequest)
				return
			}
		}
		_, err = c.DB.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
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
		if request.Date != "" {
			tmpDate, err := time.Parse("20060102", request.Date)
			if err != nil {
				utils.SendErr(w, err, http.StatusBadRequest)
				return
			}
			if request.Date == time.Now().Format("20060102") || time.Now().Before(tmpDate) {
				nextDate = request.Date
			}
		}

		res, err := c.DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
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
	case http.MethodDelete:
		id := r.FormValue("id")
		if id == "" {
			utils.SendErr(w, errors.New("id is empty"), http.StatusBadRequest)
			return
		}
		var maxID int
		row := c.DB.QueryRow(`SELECT MAX(id) FROM scheduler`)
		row.Scan(&maxID)
		if err := row.Err(); err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		givenID, err := strconv.Atoi(id)
		if err != nil {
			err = errors.New("can not parse ID")
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		if givenID > maxID {
			err = errors.New("given ID is more than number of rows in DB")
			utils.SendErr(w, err, http.StatusBadRequest)
			return
		}
		_, err = c.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	default:
		err := fmt.Errorf("no request method")
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
}
