package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetID(r)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	// TODO: Insert CheckID method here
	err = h.service.Store.CheckID(id)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	// Insert method Delete
	err = h.service.Store.Delete(id)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	h.logger.Infof("sent response via handler Task (method %s)", r.Method)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprint(w, "{}")
}
