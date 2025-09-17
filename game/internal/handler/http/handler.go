package http

import (
	"encoding/json"
	"net/http"
	"strings"

	gamectrl "ultimategaming.com/game/internal/controller/game"
)

type Handler struct {
	ctrl *gamectrl.Controller
}

func New(ctrl *gamectrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Register registra las rutas en el mux que le pases desde main.go.
func (h *Handler) Register(mux *http.ServeMux) {
	// GET /game?id=GAME_ID  -> devuelve metadata + achievements
	mux.HandleFunc("/game", h.getDetails)
}

// GET /game?id=GAME_ID
func (h *Handler) getDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	gameID := strings.TrimSpace(r.URL.Query().Get("id"))
	if gameID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	out, err := h.ctrl.GetDetails(r.Context(), gameID)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, out, http.StatusOK)
}

// Helpers

func (h *Handler) writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	// Mapeo simple de errores:
	switch {
	case err == gamectrl.ErrInvalidInput:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		// Como `game` depende de otros servicios, si algo falla al llamarlos
		// respondemos 502 Bad Gateway (puedes refinar esto si propagas códigos).
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}
}
