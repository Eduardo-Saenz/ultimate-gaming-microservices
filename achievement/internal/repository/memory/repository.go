package memory

import (
	"sync"

	"ultimategaming.com/achievement/internal/repository"
	model "ultimategaming.com/achievement/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Achievement
}

func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Achievement),
	}
}

func (r *Repository) ListByGameID(gameID string) ([]model.Achievement, error) {
    r.RLock()
    defer r.RUnlock()

    out := make([]model.Achievement, 0)
    for _, a := range r.data {
        if a.ID == gameID {
            out = append(out, *a)
        }
    }
    if len(out) == 0 {
        return nil, repository.ErrNotFound
    }
    return out, nil
}