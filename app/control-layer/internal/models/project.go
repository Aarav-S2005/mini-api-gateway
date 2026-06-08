package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission string
type LoadBalancingStrategy string

const (
	PermissionEditing Permission = "editing"
	PermissionViewing Permission = "viewing"

	RoundRobin LoadBalancingStrategy = "round-robin"
	Random     LoadBalancingStrategy = "random"
)

type Access struct {
	UserID     bson.ObjectID `bson:"user_id"`
	Permission Permission    `bson:"permission"`
}

type Project struct {
	ID                 bson.ObjectID `bson:"_id,omitempty"`
	Name               string        `bson:"name"`
	OwnerId            bson.ObjectID `bson:"owner_id"`
	GatewayApiKey      string        `bson:"gateway_api_key"`
	CreatedAt          time.Time     `bson:"created_at"`
	AccessList         []Access      `bson:"access_list"`
	Middlewares        []Middleware  `bson:"middlewares"`
	LoadBalancerConfig LoadBalancer  `bson:"load_balancer_config"`
}

type LoadBalancer struct {
	Enabled              bool                  `bson:"enabled"`
	Strategy             LoadBalancingStrategy `bson:"strategy"`
	Backends             []Backend             `bson:"backends"`
	RetryAttemptsAllowed int                   `bson:"retry_attempts_allowed"`
}

type Backend struct {
	URL    string `bson:"url"`
	Weight int    `bson:"weight"`
}
