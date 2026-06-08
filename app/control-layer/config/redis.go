package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	rdb *redis.Client
}

func NewPublisher(rdb *redis.Client) *Publisher {
	return &Publisher{
		rdb: rdb,
	}
}

func (p *Publisher) Publish(channel string, payload string) error {
	return p.rdb.Publish(context.Background(), channel, payload).Err()
}

func NewRedis(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
