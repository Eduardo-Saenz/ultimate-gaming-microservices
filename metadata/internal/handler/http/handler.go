package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"ultimategaming.com/metadata/internal/controller/metadata"
	"ultimategaming.com/metadata/internal/repository"
	"ultimategaming.com/metadata/pkg/model"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Register registra las rutas en el mux que le pases desde main.go.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/metadata", h.getOrCreate) // GET /metadata?id=GAME_ID  |  POST /metadata
}

func (h *Handler) getOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := strings.TrimSpace(r.URL.Query().Get("id"))
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		m, err := h.ctrl.Get(r.Context(), model.GameID(id))
		if err != nil {
			h.writeError(w, err)
			return
		}
		h.writeJSON(w, m, http.StatusOK)

	case http.MethodPost:
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		var in model.Metadata
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
	case errors.Is(err, metadata.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, repository.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, repository.ErrAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
