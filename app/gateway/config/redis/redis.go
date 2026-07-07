package redis

import (
	"context"
	"time"

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
