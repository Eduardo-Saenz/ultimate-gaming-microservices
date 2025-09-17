package game

import (
	"context"
	"strings"

	achievement "ultimategaming.com/game/internal/gateway/achievement/http"
	metadata "ultimategaming.com/game/internal/gateway/metadata/http"
)

// Controller usa los gateways para armar la respuesta de "detalles de juego".
type Controller struct {
	ach achievement.Client
	md  metadata.Client
}

func New(achGW achievement.Client, mdGW metadata.Client) *Controller {
	return &Controller{ach: achGW, md: mdGW}
}

// GameDetails es el DTO de salida del servicio "game".
type GameDetails struct {
	Metadata     *metadata.MetadataDTO         `json:"metadata"`
	Achievements []achievement.AchievementDTO  `json:"achievements"`
}

// GetDetails regresa metadata + achievements de un juego.
func (c *Controller) GetDetails(ctx context.Context, gameID string) (*GameDetails, error) {
	if strings.TrimSpace(gameID) == "" {
		return nil, ErrInvalidInput
	}

	md, err := c.md.Get(gameID)
	if err != nil {
		return nil, err
	}

	achs, err := c.ach.ListByGame(gameID)
	if err != nil {
		return nil, err
	}

	return &GameDetails{
		Metadata:     md,
		Achievements: achs,
	}, nil
}
