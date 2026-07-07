package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LoadBalancingStrategy string

const (
	RoundRobinLoadBalancing LoadBalancingStrategy = "ROUND_ROBIN"
	RandomLoadBalancing     LoadBalancingStrategy = "RANDOM"
	IPHashLoadBalancing     LoadBalancingStrategy = "IP_HASH"
	WeightedRoundRobin      LoadBalancingStrategy = "WEIGHT_ROUND_ROBIN"
)

type Service struct {
	ID                    bson.ObjectID         `bson:"_id,omitempty"`
	ProjectID             bson.ObjectID         `bson:"project_id"`
	Name                  string                `bson:"name"`
	LoadBalancingStrategy LoadBalancingStrategy `bson:"strategy"`
	Backends              []Backend             `bson:"backends"`
	CreatedAt             time.Time             `bson:"created_at"`
}

type Backend struct {
	URL    string `bson:"url"`
	Weight int    `bson:"weight,omitempty"`
}
