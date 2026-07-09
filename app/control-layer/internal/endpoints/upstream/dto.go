package route

import (
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateOrUpdateUpstreamRequestDTO struct {
	Name                  string                       `json:"name"`
	LoadBalancingStrategy models.LoadBalancingStrategy `json:"load_balancing_strategy"`
	Backends              []models.Backend             `json:"backends"`
}

type GetUpstreamResponseDTO struct {
	ID                    bson.ObjectID                `json:"_id"`
	Name                  string                       `json:"name"`
	LoadBalancingStrategy models.LoadBalancingStrategy `json:"load_balancing_strategy"`
	Backends              []models.Backend             `json:"backends"`
}

type GetAllUpstreamResponseDTO struct {
	upstreams []GetUpstreamResponseDTO
}

type CreateUpstreamResponseDTO struct {
	UpstreamID bson.ObjectID `json:"project_id"`
}
