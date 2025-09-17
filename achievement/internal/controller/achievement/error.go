package achievement

import "errors"

// Errores propios de la capa de aplicación (validaciones/reglas)
var (
	ErrInvalidInput = errors.New("invalid input")
)
