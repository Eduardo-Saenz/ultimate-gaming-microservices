package repository

import "errors"

// Errores de dominio del repositorio
var (
    ErrNotFound      = errors.New("not found")
    ErrAlreadyExists = errors.New("already exists")
)
