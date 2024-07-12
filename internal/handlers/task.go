package handlers

import (
	"fmt"
	"net/http"

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
func (h *Handler) Task(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTask(w, r)
	case http.MethodPut:
		h.EditTask(w, r)
	case http.MethodPost:
		h.AddTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		err := fmt.Errorf("method not allowed")
		h.SendErr(w, err, http.StatusMethodNotAllowed)
		return
	}
}
