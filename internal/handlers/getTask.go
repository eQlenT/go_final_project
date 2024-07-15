package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetID(r)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	task, err := h.service.GetTask(id)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(task)
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
