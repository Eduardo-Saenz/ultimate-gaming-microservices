package repository

import (
	"context"

	"ultimategaming.com/metadata/pkg/model"
)

// Repository define el acceso a datos para Metadata.
type Repository interface {
	// Get devuelve la metadata de un juego por ID.
	Get(ctx context.Context, id model.GameID) (*model.Metadata, error)

	// Create inserta una nueva metadata (ID único).
	Create(ctx context.Context, m model.Metadata) error
}
