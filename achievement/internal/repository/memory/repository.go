package memory

import (
	"context"
	"sync"

	"ultimategaming.com/achievement/internal/repository"
	"ultimategaming.com/achievement/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[model.AchievementID]*model.Achievement
}

func New() *Repository {
	return &Repository{
		data: make(map[model.AchievementID]*model.Achievement),
	}
}

// Get devuelve un logro por su ID.
func (r *Repository) Get(ctx context.Context, id model.AchievementID) (*model.Achievement, error) {
	r.RLock()
	defer r.RUnlock()

	a, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	// devolvemos una copia para evitar mutaciones externas del mapa interno
	cp := *a
	return &cp, nil
}

// ListByGame devuelve todos los logros asociados a un GameID.
func (r *Repository) ListByGame(ctx context.Context, gameID model.GameID) ([]model.Achievement, error) {
	r.RLock()
	defer r.RUnlock()

	out := make([]model.Achievement, 0)
	for _, a := range r.data {
		if a.GameID == gameID {
			out = append(out, *a) // copiamos valor
		}
	}

	if len(out) == 0 {
		return nil, repository.ErrNotFound
	}
	return out, nil
}

// Create inserta un nuevo logro si no existe el ID.
func (r *Repository) Create(ctx context.Context, a model.Achievement) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.data[a.ID]; exists {
		return repository.ErrAlreadyExists
	}

	// guardamos copia para aislar memoria interna
	cp := a
	r.data[a.ID] = &cp
	return nil
}
