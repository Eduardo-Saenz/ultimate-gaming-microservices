package achievement

import (
	"errors"
	"strings"

	"ultimategaming.com/achievement/pkg/model"
)

// Repository define lo que el Controller necesita del almacenamiento.
// La implementación concreta (memoria, BD, etc.) vendrá después.
type Repository interface {
	ListByGameID(gameID string) ([]model.Achievement, error)
	Add(gameID string, a model.Achievement) error
}

type Controller struct {
	repo Repository
}

func New(r Repository) *Controller {
	return &Controller{repo: r}
}

// List devuelve los logros de un juego.
func (c *Controller) List(gameID string) ([]model.Achievement, error) {
	gameID = strings.TrimSpace(gameID)
	if gameID == "" {
		return nil, errors.New("game_id is required")
	}
	return c.repo.ListByGameID(gameID)
}

// Add agrega un logro al juego (validación mínima).
func (c *Controller) Add(gameID string, a model.Achievement) error {
	gameID = strings.TrimSpace(gameID)
	if gameID == "" {
		return errors.New("game_id is required")
	}
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.Name) == "" {
		return errors.New("achievement id and name are required")
	}
	return c.repo.Add(gameID, a)
}
