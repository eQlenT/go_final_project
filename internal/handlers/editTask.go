package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go_final_project/internal/models/service/store/task"
)

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	var task task.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		err = fmt.Errorf("can't parse response")
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	err = task.CheckTask()
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	} else {
		id, err := strconv.Atoi(task.ID)
		if err != nil {
			err = fmt.Errorf("cat not parse id")
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		// TODO: CHANGE TO SERVICE METHOD
		err = h.service.Store.CheckID(id)
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		task.Date, err = task.CheckDate()
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
	}
	err = h.service.Update(&task)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	h.logger.Infof("sent response via handler Task (method %s)", r.Method)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write([]byte("{}"))
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
}
