package cachedb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type cacheDB struct {
	Master  *redis.Client
	Replica *redis.Client
}

var cachedb *cacheDB

type Config struct {
	MasterHost      string
	MasterPassword  string
	MasterDB        int
	ReplicaHost     string
	ReplicaPassword string
	ReplicaDB       int
}

func Open(cfg Config) (*cacheDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rdbm := redis.NewClient(&redis.Options{
		Addr:     cfg.MasterHost,
		Password: cfg.MasterPassword,
		DB:       cfg.MasterDB,
	})
	if err := rdbm.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	rdbr := redis.NewClient(&redis.Options{
		Addr:     cfg.ReplicaHost,
		Password: cfg.ReplicaPassword,
		DB:       cfg.ReplicaDB,
	})
	if err := rdbr.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	cachedb = &cacheDB{
		Master:  rdbm,
		Replica: rdbr,
	}
	return cachedb, nil
}

// Store manages the set of APIs for user database access.
type Store struct {
	log     *zap.SugaredLogger
	cachedb cacheDB
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, cachedb cacheDB) *Store {
	return &Store{
		log:     log,
		cachedb: cachedb,
	}
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := cachedb.Master.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("set session token: %w", err)
	}

	return nil
}

func Get(ctx context.Context, key string) (string, error) {
	val, err := cachedb.Replica.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("get session token: %w", err)
	}

	return val, nil
}

func XRange(ctx context.Context, streamName string, start string, stop string) ([]redis.XMessage, error) {
	val, err := cachedb.Replica.XRange(ctx, streamName, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("get session token: %w", err)
	}

	return val, nil
}

func Subscribe(ctx context.Context, channelName string) *redis.PubSub {
	return cachedb.Replica.Subscribe(ctx, channelName)
}

func XAdd(ctx context.Context, streamName string, values map[string]interface{}) (string, error) {
	val, err := cachedb.Master.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: values,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("xadd: %w", err)
	}

	return val, nil
}

func Publish(ctx context.Context, channelName string, message string) error {
	err := cachedb.Master.Publish(ctx, channelName, message).Err()
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
