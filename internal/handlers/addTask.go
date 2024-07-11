package handlers

import (
	"encoding/json"
	"go_final_project/internal/models"
	"net/http"
)

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var id struct {
		ID int `json:"id"`
	}
	var request models.Task
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

	nextDate, err := request.CompleteRequest()
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	if request.Date != "" {
		nextDate, err = request.CheckDate()
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
		}
	}
	request.Date = nextDate
	lastInsertID, err := h.service.Store.Insert(&request)
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
	w.Write(response)
}
