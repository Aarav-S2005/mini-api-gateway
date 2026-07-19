package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const ConfigUpdatesChannel = "gateway:config:updates"

type ResourceType string
type Operation string

const (
	ResourceProject  ResourceType = "project"
	ResourceRoute    ResourceType = "route"
	ResourceUpstream ResourceType = "upstream"

	WriteOperation   Operation = "write"
	DeleteOperation  Operation = "delete"
	UpdateOperation  Operation = "update"
)

type Publisher struct {
	client *redis.Client
}

func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{
		client: client,
	}
}

func (p *Publisher) Publish(ctx context.Context, payload UpdateEventNotification) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.client.Publish(ctx, ConfigUpdatesChannel, data).Err()
}

type UpdateEventNotification struct {
	Resource         ResourceType  `json:"resource"`
	Operation  Operation
	ResourceID       bson.ObjectID `json:"resource_id"`
	ConfigUpdateTime time.Time     `json:"config_update_time"`
}
