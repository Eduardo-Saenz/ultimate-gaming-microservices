package repository

import (
	"context"

	"ultimategaming.com/achievement/pkg/model"
)

// Repository define qué necesita el dominio para acceder a datos.
type Repository interface {
    Get(ctx context.Context, id model.AchievementID) (*model.Achievement, error)
    ListByGame(ctx context.Context, gameID model.GameID) ([]model.Achievement, error)
    Create(ctx context.Context, a model.Achievement) error
}
