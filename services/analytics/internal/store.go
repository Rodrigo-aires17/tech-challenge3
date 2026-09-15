package internal

import (
	"context"
	"sync"
	"time"
)

// Event é o mesmo payload publicado pelo serviço `evaluation` na fila SQS.
type Event struct {
	FlagID    string    `json:"flag_id"`
	UserID    string    `json:"user_id"`
	Matched   bool      `json:"matched"`
	Timestamp time.Time `json:"timestamp"`
}

// Store abstrai a agregação/persistência das métricas. Em produção,
// DynamoDBStore grava na tabela ToggleMasterAnalytics; em testes usamos
// InMemoryStore.
type Store interface {
	Record(ctx context.Context, event Event) error
	Counts(ctx context.Context, flagID string) (matched, total int, err error)
}

type flagCounter struct {
	matched int
	total   int
}

type InMemoryStore struct {
	mu       sync.Mutex
	counters map[string]*flagCounter
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{counters: make(map[string]*flagCounter)}
}

func (s *InMemoryStore) Record(_ context.Context, event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.counters[event.FlagID]
	if !ok {
		c = &flagCounter{}
		s.counters[event.FlagID] = c
	}
	c.total++
	if event.Matched {
		c.matched++
	}
	return nil
}

func (s *InMemoryStore) Counts(_ context.Context, flagID string) (int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.counters[flagID]
	if !ok {
		return 0, 0, nil
	}
	return c.matched, c.total, nil
}
