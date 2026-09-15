package internal

import "context"

// Cache abstrai o armazenamento do resultado de avaliação, evitando recalcular
// o hash de rollout a cada requisição. Em produção, RedisCache implementa esta
// interface contra o ElastiCache; em testes usamos InMemoryCache.
type Cache interface {
	Get(ctx context.Context, key string) (bool, bool, error) // valor, encontrado, erro
	Set(ctx context.Context, key string, value bool) error
}

type InMemoryCache struct {
	data map[string]bool
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{data: make(map[string]bool)}
}

func (c *InMemoryCache) Get(_ context.Context, key string) (bool, bool, error) {
	v, ok := c.data[key]
	return v, ok, nil
}

func (c *InMemoryCache) Set(_ context.Context, key string, value bool) error {
	c.data[key] = value
	return nil
}
