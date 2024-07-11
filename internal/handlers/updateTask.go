package handlers

import (
	"encoding/json"
	"fmt"
	"go_final_project/internal/models"
	"net/http"
)

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	err = task.CheckTask()
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	} else {
		err = h.service.Store.CheckID(task.ID)
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
		}
	}
	err = h.service.Store.Update(&task)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
	}
	h.logger.Infof("sent response via handler Task (method %s)", r.Method)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprint(w, "{}")
}
