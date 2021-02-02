package blacklist

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type BlackList interface {
	Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	CheckIfExist(ctx context.Context, key string) error
}

type blackList struct {
	client *redis.Client
}

// Add new refresh token to db
func (b blackList) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return b.client.SetNX(ctx, key, value, ttl).Err()
}

// CheckIfExist if refresh token existing in db. If refresh token not existed return error
func (b blackList) CheckIfExist(ctx context.Context, key string) error {
	_, err := b.client.Get(ctx, key).Result()
	switch err {
	case redis.Nil:
		return fmt.Errorf("%s", "key does not exists")
	case nil:
		return nil
	default:
		return err
	}
}

func New(client *redis.Client) *blackList {
	return &blackList{
		client: client,
	}
}
