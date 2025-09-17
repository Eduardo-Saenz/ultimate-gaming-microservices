package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"ultimategaming.com/achievement/internal/controller/achievement"
	"ultimategaming.com/achievement/internal/repository"
	"ultimategaming.com/achievement/pkg/model"
)

type Handler struct {
	ctrl *achievement.Controller
}

func New(ctrl *achievement.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Register registra las rutas en el mux que le pases desde main.go
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/achievement", h.getAchievement)      // GET /achievement?id=ACHV_ID
	mux.HandleFunc("/achievements", h.listOrCreate)       // GET /achievements?gameId=GAME_ID  |  POST /achievements
}

// GET /achievement?id=ACHV_ID
func (h *Handler) getAchievement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	a, err := h.ctrl.Get(r.Context(), model.AchievementID(id))
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, a, http.StatusOK)
}

// GET /achievements?gameId=GAME_ID
// POST /achievements   (JSON body)
func (h *Handler) listOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		gameID := strings.TrimSpace(r.URL.Query().Get("gameId"))
		if gameID == "" {
			http.Error(w, "missing gameId", http.StatusBadRequest)
			return
		}
		list, err := h.ctrl.ListByGame(r.Context(), model.GameID(gameID))
		if err != nil {
			h.writeError(w, err)
			return
		}
		h.writeJSON(w, list, http.StatusOK)

	case http.MethodPost:
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		var in model.Achievement
		if err := json.Unmarshal(body, &in); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := h.ctrl.Create(r.Context(), in); err != nil {
			h.writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Helpers

func (h *Handler) writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, achievement.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, repository.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
