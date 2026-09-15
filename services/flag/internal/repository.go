package internal

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("feature flag not found")

type Flag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Repository abstrai a persistência das flags. Em produção, PostgresRepository
// (flag_db) implementa esta interface; nos testes usamos InMemoryRepository.
type Repository interface {
	Create(ctx context.Context, name, description string, enabled bool) (Flag, error)
	Get(ctx context.Context, id string) (Flag, error)
	List(ctx context.Context) ([]Flag, error)
	Update(ctx context.Context, id string, enabled bool, description string) (Flag, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	flags map[string]Flag
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{flags: make(map[string]Flag)}
}

func (r *InMemoryRepository) Create(_ context.Context, name, description string, enabled bool) (Flag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	flag := Flag{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Enabled:     enabled,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.flags[flag.ID] = flag
	return flag, nil
}

func (r *InMemoryRepository) Get(_ context.Context, id string) (Flag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, ok := r.flags[id]
	if !ok {
		return Flag{}, ErrNotFound
	}
	return flag, nil
}

func (r *InMemoryRepository) List(_ context.Context) ([]Flag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Flag, 0, len(r.flags))
	for _, f := range r.flags {
		result = append(result, f)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(_ context.Context, id string, enabled bool, description string) (Flag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	flag, ok := r.flags[id]
	if !ok {
		return Flag{}, ErrNotFound
	}
	flag.Enabled = enabled
	flag.Description = description
	flag.UpdatedAt = time.Now().UTC()
	r.flags[id] = flag
	return flag, nil
}

func (r *InMemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.flags[id]; !ok {
		return ErrNotFound
	}
	delete(r.flags, id)
	return nil
}
