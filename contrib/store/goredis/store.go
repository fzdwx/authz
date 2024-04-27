package goredis

import (
	"context"
	"errors"
	"fmt"
	"github.com/fzdwx/authz"
	"github.com/redis/go-redis/v9"
)

type store struct {
	client *redis.Client
}

func (s *store) Get(ctx context.Context, key string) (string, error) {
	if val, err := s.client.Get(ctx, key).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", authz.ErrStoreValueNotFound
		}
		return "", fmt.Errorf("goredis get key %s: %w", key, err)
	} else {
		return val, nil
	}
}

func (s *store) Set(ctx context.Context, key string, value string) error {
	if err := s.client.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("goredis set key %s: %w", key, err)
	}
	return nil
}

func NewStore(client *redis.Client) authz.Store {
	return &store{
		client: client,
	}
}
