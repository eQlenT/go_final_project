package handlers

import (
	"errors"
	"fmt"
	"net/http"
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
func (h *Handler) TaskDone(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if r.Method != http.MethodPost {
		err := errors.New("request method must be post")
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	id, err := h.GetID(r)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
	}
	h.service.Done(id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprint(w, "{}")
}
