package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ultimategaming.com/achievement/internal/controller/achievement"
	"ultimategaming.com/achievement/internal/repository"
)

type Handler struct {
	ctrl *achievement.Controller
}

func New(ctrl *achievement.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetAchievement(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	a, err := h.ctrl.List(id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(a); err != nil {
		log.Printf("Response error: %v\n", err)
	}
}
