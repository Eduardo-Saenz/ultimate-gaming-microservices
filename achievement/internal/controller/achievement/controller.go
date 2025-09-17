package achievement

import (
	"context"
	"strings"

	"ultimategaming.com/achievement/internal/repository"
	"ultimategaming.com/achievement/pkg/model"
)

type Controller struct {
	repo repository.Repository
}

func New(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

// Get devuelve un logro por su ID.
func (c *Controller) Get(ctx context.Context, id model.AchievementID) (*model.Achievement, error) {
	if strings.TrimSpace(string(id)) == "" {
		return nil, ErrInvalidInput
	}
	return c.repo.Get(ctx, id)
}

// ListByGame devuelve todos los logros de un GameID.
func (c *Controller) ListByGame(ctx context.Context, gameID model.GameID) ([]model.Achievement, error) {
	if strings.TrimSpace(string(gameID)) == "" {
		return nil, ErrInvalidInput
	}
	return c.repo.ListByGame(ctx, gameID)
}

// Create valida y crea un logro.
func (c *Controller) Create(ctx context.Context, a model.Achievement) error {
	if strings.TrimSpace(string(a.ID)) == "" ||
		strings.TrimSpace(string(a.GameID)) == "" ||
		strings.TrimSpace(a.Name) == "" {
		return ErrInvalidInput
	}
	if a.Points < 0 {
		return ErrInvalidInput
	}
	return c.repo.Create(ctx, a)
}
