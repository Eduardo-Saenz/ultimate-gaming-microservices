package achievement

import (
	"errors"
	"strings"

	"ultimategaming.com/achievement/pkg/model"
)

type Repository interface {
	ListByGameID(gameID string) ([]model.Achievement, error)
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