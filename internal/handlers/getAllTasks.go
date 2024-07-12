// Реализовать обработчик для GET-запроса /api/tasks.
// Он должен возвращать список ближайших задач в формате JSON в виде списка в поле tasks.
// Задачи должны быть отсортированы по дате в сторону увеличения.
// Каждая задача должна содержать все поля таблицы scheduler в виде строк.
// Дата представлена в уже знакомом вам формате 20060102.

package handlers

import (
	"encoding/json"
	"net/http"

	_ "modernc.org/sqlite"

	"go_final_project/internal/models/service/store/task"
)

// GetTasks - обработчик для GET-запросов к /api/tasks.
// Он извлекает список ближайших задач из базы данных и возвращает их в виде JSON-ответа.
// Задачи сортируются по дате в порядке возрастания.
// Каждая задача содержит все поля таблицы scheduler в виде строк.
// Дата представлена в формате 20060102.
//
// Параметры:
// - w: http.ResponseWriter для записи ответа.
// - r: *http.Request, содержащий данные запроса.
//
// Возвращает:
// // - Функция не возвращает никакого значения, но записывает JSON-ответ в http.ResponseWriter.
func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	const limit = 50
	var tasks map[string][]task.Task
	search := r.FormValue("search")
	var isSearch bool = search != ""
	var err error
	if isSearch {
		tasks, err = h.service.Search(search, limit)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		tasks, err = h.service.Store.GetAll(limit)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}
	}
	response, err := json.Marshal(tasks)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	h.logger.Infof("sent response via handler GetTasks")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(response)
	if err != nil {
		h.logger.Error(err)
	}
}
