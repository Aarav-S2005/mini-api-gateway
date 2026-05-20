package redis

import (
	"context"

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
