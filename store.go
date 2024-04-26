package authz

import (
	"context"
	"errors"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type memoryStore struct {
	data map[string]string
}

func (m *memoryStore) Set(ctx context.Context, key string, value string) error {
	m.data[key] = value
	return nil
}

var (
	ErrStoreValueNotFound = errors.New("value not found")
)

func (m *memoryStore) Get(ctx context.Context, key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", ErrStoreValueNotFound
	}
	return val, nil
}

func NewMemoryStore() Store {
	return &memoryStore{
		data: make(map[string]string),
	}
}
