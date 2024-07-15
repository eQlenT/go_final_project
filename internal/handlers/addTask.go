package handlers

import (
	"encoding/json"
	"net/http"

	"go_final_project/internal/models/service/store/task"
)

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var id struct {
		ID int `json:"id"`
	}
	var request task.Task
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}

	err = request.CheckTask()
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	// nextDate, err := request.CompleteRequest()
	// if err != nil {
	// 	h.SendErr(w, err, http.StatusBadRequest)
	// 	return
	// }
	// request.Date = nextDate
	// if request.Date != "" {
	// 	nextDate, err = request.CheckDate()
	// 	if err != nil {
	// 		h.SendErr(w, err, http.StatusBadRequest)
	// 	}
	// 	request.Date = nextDate
	// }
	err = h.service.DateToAdd(&request)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	lastInsertID, err := h.service.Insert(&request)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	id.ID = lastInsertID
	response, err := json.Marshal(id)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	h.logger.Infof("sent response via handler Task (method %s)", r.Method)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(response)
	if err != nil {
		h.logger.Error(err)
	}
}
