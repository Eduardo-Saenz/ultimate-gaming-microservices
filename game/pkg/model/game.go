package model

import (
	achievementModel "ultimategaming.com/achievement/pkg/model"
	metadataModel "ultimategaming.com/metadata/pkg/model"
)

type GameDetails struct {
	Metadata     metadataModel.Metadata         `json:"metadata"`
	Achievements []achievementModel.Achievement `json:"achievements,omitempty"`
}
