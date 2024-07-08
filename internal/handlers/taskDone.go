package handlers

import (
	"errors"
	"fmt"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TaskDone обрабатывает завершение задачи. Если задача повторяется, она обновляет дату для следующего повторения.
// В противном случае, она удаляет задачу из планировщика.
//
// Параметры:
// - w: http.ResponseWriter для записи ответа.
// - r: http.Request, содержащий данные запроса.
//
// Возвращает:
// - Записывает JSON-ответ с пустым объектом в ответный writer.
// - Если во время процесса возникает ошибка, он отправляет ответ с ошибкой с соответствующим кодом состояния.
func (c *DBConnection) TaskDone(w http.ResponseWriter, r *http.Request) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if r.Method != http.MethodPost {
		err := errors.New("request method must be post")
		utils.SendErr(w, err, http.StatusInternalServerError)
		return
	}
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
	if task.Repeat != "" {
		nextDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		if task.Date == time.Now().Format("20060102") {
			date, err := time.Parse("20060102", task.Date)
			if err != nil {
				utils.SendErr(w, err, http.StatusInternalServerError)
				return
			}
			rptSlc := strings.Split(task.Repeat, " ")
			subDays, err := strconv.Atoi(rptSlc[1])
			nextDate = date.AddDate(0, 0, subDays).Format("20060102")
			if err != nil {
				utils.SendErr(w, err, http.StatusInternalServerError)
				return
			}
		}
		_, err = c.DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`,
			nextDate, task.ID)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	} else {
		_, err = c.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, task.ID)
		if err != nil {
			utils.SendErr(w, err, http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	}
}
