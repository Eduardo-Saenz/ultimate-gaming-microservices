package model

type AchievementID string
type GameID string

type Achievement struct {
	ID          AchievementID `json:"id"`          // Identificador único del logro
	GameID      GameID        `json:"gameId"`      // Identificador del juego al que pertenece
	Name        string        `json:"name"`        // Nombre del logro
	Description string        `json:"description"` // Descripción del logro
	Points      int           `json:"points"`      // Puntos que otorga
	Secret      bool          `json:"secret"`      // Indica si es secreto/especial
}
