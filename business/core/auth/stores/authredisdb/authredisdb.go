// Package authredisdb contains redis related Set and Get functionality.
package authredisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisDB struct {
	Master  *redis.Client
	Replica *redis.Client
}

// Store manages the set of APIs for user database access.
type Store struct {
	log     *zap.SugaredLogger
	redisdb redisDB
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, redisdb redisDB) *Store {
	return &Store{
		log:     log,
		redisdb: redisdb,
	}
}

func (s *Store) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := s.redisdb.Master.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("set session token: %w", err)
	}

	return nil
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	val, err := s.redisdb.Replica.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("get session token: %w", err)
	}

	return val, nil
}
