package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Publisher struct {
	client *redis.Client
}

func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{
		client: client,
	}
}

func (p *Publisher) Publish(ctx context.Context, channel string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.client.Publish(ctx, channel, data).Err()
}

type UpdateEventNotification struct {
	Resource                 string
	ResourceID               bson.ObjectID
	ResourceAttributeUpdated []string
	ConfigUpdateTime         time.Time
}
