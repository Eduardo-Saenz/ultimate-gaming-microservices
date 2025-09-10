package game

import (
	"context"
	"errors"

	achievementModel "ultimategaming.com/achievement/pkg/model"
	"ultimategaming.com/game/internal/gateway"
	"ultimategaming.com/game/pkg/model"
	metadataModel "ultimategaming.com/metadata/pkg/model"
)

var ErrNotFound = errors.New("Game metadata not found")

type achievementGateway interface {
	// Lista los logros de un record (juego).
	ListByRecord(ctx context.Context, recordID achievementModel.RecordID, recordType achievementModel.RecordType) ([]achievementModel.Achievement, error)
}

type metadataGateway interface {
	// Obtiene la metadata de un juego por ID.
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

type Controller struct {
	achievementGateway achievementGateway
	metadataGateway    metadataGateway
}

func New(achievementGateway achievementGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{achievementGateway, metadataGateway}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.GameDetails, error) {
	// 1) Traer metadata del juego
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &model.GameDetails{Metadata: *metadata}

	achievements, err := c.achievementGateway.ListByRecord(
		ctx,
		achievementModel.RecordID(id),
		achievementModel.RecordTypeGame,
	)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		return nil, err
	} else if err == nil {
		details.Achievements = achievements
	}

	return details, nil

}
