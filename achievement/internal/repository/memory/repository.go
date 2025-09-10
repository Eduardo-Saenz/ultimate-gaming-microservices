package memory

import (
	"context"
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

func (r *Repository) Get(_ context.Context, id string) (*model.Achievement, error) {
	r.RLock()
	defer r.RUnlock()

	a, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return a, nil
}

func (r *Repository) Put(_ context.Context, id string, a *model.Achievement) error {
	r.Lock()
	defer r.Unlock()

	r.data[id] = a
	return nil
}
