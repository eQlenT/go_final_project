package handlers

import (
	"fmt"
	"go_final_project/internal/models"
	"net/http"
	"strconv"
	"sync"

	"go.uber.org/zap"
)

type Handler struct {
	service *models.TaskService
	logger  *zap.SugaredLogger
	mu      sync.Mutex
}

func NewHandler(service *models.TaskService, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		mu:      sync.Mutex{},
	}
}

// SendErr - это метод для типа Handler, которая отправляет клиенту ответ об ошибке.
// Она устанавливает соответствующий HTTP-код состояния и тип содержимого, а также форматирует сообщение об ошибке как JSON-объект.
//
// Параметры:
// w - http.ResponseWriter, куда будет записан ответ.
// err - объект ошибки, содержащий сообщение об ошибке для отправки.
// status - целое число, представляющее HTTP-код состояния для отправки.
//
// Возвращает:
// Эта функция не возвращает никакого значения.
func (h *Handler) SendErr(w http.ResponseWriter, err error, status int) {
	h.logger.Error(err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), status)
}

func (h *Handler) GetID(r *http.Request) (int, error) {
	var id int
	var err error
	idStr := r.FormValue("id")
	if idStr == "" {
		return 0, fmt.Errorf("id is empty")
	}
	if idStr != "" || len(idStr) != 0 {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			return 0, fmt.Errorf("can not parse ID")
		}
	}
	return id, nil
}
