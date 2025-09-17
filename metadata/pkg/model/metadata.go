package model

// Tipo fuerte para el ID del juego
type GameID string

// Metadata describe la información principal de un videojuego.
type Metadata struct {
	ID          GameID `json:"id"`                    // Identificador único del juego
	Name        string `json:"name"`                  // Nombre del juego
	Description string `json:"description"`           // Descripción corta
	Genre       string `json:"genre,omitempty"`       // Género principal (opcional)
	Developer   string `json:"developer,omitempty"`   // Desarrollador (opcional)
	ReleaseYear int    `json:"releaseYear,omitempty"` // Año de lanzamiento (opcional)
}
