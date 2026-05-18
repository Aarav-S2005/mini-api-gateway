package redis_cfg

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	rdb *redis.Client
}

func NewSubscriber(rdb *redis.Client) *Subscriber {
	return &Subscriber{rdb: rdb}
}

func (subscriber *Subscriber) Subscribe(ctx context.Context, channel string, handler func(msg []byte)) {
	sub := subscriber.rdb.Subscribe(ctx, channel)

	ch := sub.Channel()
	for msg := range ch {
		handler([]byte(msg.Payload))
	}
}
