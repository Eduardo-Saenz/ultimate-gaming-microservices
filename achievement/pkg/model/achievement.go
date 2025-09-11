package model

type ID string
type RecordID string
type RecordType string

const (
	RecordTypeGame = RecordType("game")
)

// Achievement representa un logro disponible en un juego.
type Achievement struct {
	ID          string     `json:"game_id"`
	Name        string `json:"name"`
	Points      int    `json:"points"`
	Description string `json:"description"`
	Secret      bool   `json:"secret,omitempty"`
	Icon        string `json:"icon,omitempty"`
}
