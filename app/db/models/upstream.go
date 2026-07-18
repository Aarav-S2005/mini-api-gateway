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
	WeightedRoundRobin      LoadBalancingStrategy = "WEIGHTED_ROUND_ROBIN"
	LeastConnections        LoadBalancingStrategy = "LEAST_CONNECTIONS"
)

type Upstream struct {
	ID                    bson.ObjectID         `bson:"_id,omitempty"`
	ProjectID             bson.ObjectID         `bson:"project_id"`
	Name                  string                `bson:"name"`
	LoadBalancingStrategy LoadBalancingStrategy `bson:"load_balancing_strategy"`
	Backends              []Backend             `bson:"backends"`
	CreatedAt             time.Time             `bson:"created_at"`
	UpdatedAt             time.Time             `bson:"updated_at"`
}

type Backend struct {
	URL    string `bson:"url" json:"url"`
	Weight *int   `bson:"weight,omitempty" json:"weight,omitempty"`
}
