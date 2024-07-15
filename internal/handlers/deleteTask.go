package handlers

import "net/http"

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetID(r)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	err = h.service.CheckID(id)
	if err != nil {
		h.SendErr(w, err, http.StatusBadRequest)
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
	h.logger.Infof("sent response via handler Task (method %s)", r.Method)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write([]byte("{}"))
	if err != nil {
		h.logger.Error(err)
	}
}
