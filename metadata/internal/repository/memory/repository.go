package memory

import (
	"context"
	"sync"

	"ultimategaming.com/metadata/internal/repository"
	"ultimategaming.com/metadata/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[model.GameID]*model.Metadata
}

func New() *Repository {
	return &Repository{
		data: make(map[model.GameID]*model.Metadata),
	}
}

func (r *Repository) Get(ctx context.Context, id model.GameID) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	cp := *m
	return &cp, nil
}

func (r *Repository) Create(ctx context.Context, m model.Metadata) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.data[m.ID]; exists {
		return repository.ErrAlreadyExists
	}
	cp := m
	r.data[m.ID] = &cp
	return nil
}
