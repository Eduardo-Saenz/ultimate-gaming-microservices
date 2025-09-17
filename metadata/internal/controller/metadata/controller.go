package metadata

import (
	"context"
	"strings"

	"ultimategaming.com/metadata/internal/repository"
	"ultimategaming.com/metadata/pkg/model"
)

type Controller struct {
	repo repository.Repository
}

func New(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

// Get devuelve la metadata de un juego por su ID.
func (c *Controller) Get(ctx context.Context, id model.GameID) (*model.Metadata, error) {
	if strings.TrimSpace(string(id)) == "" {
		return nil, ErrInvalidInput
	}
	return c.repo.Get(ctx, id)
}

// Create valida e inserta una nueva metadata.
func (c *Controller) Create(ctx context.Context, m model.Metadata) error {
	if strings.TrimSpace(string(m.ID)) == "" ||
		strings.TrimSpace(m.Name) == "" {
		return ErrInvalidInput
	}
	return c.repo.Create(ctx, m)
}
