package config

import (
	"context"
	"encoding/json"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/redis/go-redis/v9"
)

const ConfigUpdatesChannel = "gateway:config:updates"

type ResourceType string

const (
	ResourceProject  ResourceType = "project"
	ResourceRoute    ResourceType = "route"
	ResourceUpstream ResourceType = "upstream"
)

type Subscriber struct {
	rdb *redis.Client
}

func NewSubscriber(rdb *redis.Client) *Subscriber {
	return &Subscriber{rdb: rdb}
}

func (subscriber *Subscriber) Subscribe(ctx context.Context, handler func(context.Context, config.UpdateEventNotification) error,
) error {
	sub := subscriber.rdb.Subscribe(ctx, ConfigUpdatesChannel)
	defer sub.Close()

	if _, err := sub.Receive(ctx); err != nil {
		return err
	}
	ch := sub.Channel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			var notification config.UpdateEventNotification
			if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
				continue
			}
			if err := handler(ctx, notification); err != nil {
				return err
			}
		}
	}
}
