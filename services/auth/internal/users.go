package internal

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// UserStore mantém usuários em memória com senhas sempre em hash (nunca texto plano).
// Em produção, isto seria substituído por um repositório PostgreSQL (auth_db).
type UserStore struct {
	mu    sync.RWMutex
	users map[string]string // username -> bcrypt hash
}

func NewUserStore() *UserStore {
	return &UserStore{users: make(map[string]string)}
}

func (s *UserStore) Register(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[username] = string(hash)
	return nil
}

func (s *UserStore) Authenticate(username, password string) error {
	s.mu.RLock()
	hash, ok := s.users[username]
	s.mu.RUnlock()

	if !ok {
		return ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}
