package route

import (
	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateOrUpdateRouteRequestDTO struct {
	Path         string          `json:"path"`
	PathType     models.PathType `json:"path_type"`
	Method       string          `json:"method"`
	UpstreamName string          `json:"upstream_name"`
	AuthMode     models.AuthMode `json:"auth_mode"`
	Enabled      bool            `json:"enabled"`
}

type CreateRouteResponseDTO struct {
	RouteID bson.ObjectID `json:"route_id"`
}

type GetRouteResponseDTO struct {
	ID           bson.ObjectID   `json:"route_id"`
	Path         string          `json:"path"`
	PathType     models.PathType `json:"path_type"`
	Method       string          `json:"method"`
	UpstreamName string          `json:"upstream_name"`
	AuthMode     models.AuthMode `json:"auth_mode"`
	Enabled      bool            `json:"enabled"`
}

type GetAllRoutesResponseDTO struct {
	Routes []GetRouteResponseDTO `json:"routes"`
}
