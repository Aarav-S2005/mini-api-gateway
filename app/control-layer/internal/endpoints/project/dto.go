package project

import (
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreatProjectRequest struct {
	Name       string          `json:"name"`
	AccessList []models.Access `json:"access_list"`
}

type CreatProjectResponse struct {
	Id       bson.ObjectID `json:"_id"`
	ApiGwKey string        `json:"api_gw_key"`
}

type UpdateMiddlewaresRequest struct {
	Middlewares []models.Middleware `json:"middleware"`
}

type UpdateAccessListRequest struct {
	AccessList []models.Access `json:"access_list"`
}

type GetProjectResponse struct {
	ID          bson.ObjectID       `json:"_id,omitempty"`
	Name        string              `json:"name"`
	CreatedAt   time.Time           `json:"created_at"`
	Middlewares []models.Middleware `json:"middlewares"`
	Permission  string              `json:"permission"`
}

type ListProjectResponse struct {
	ID   bson.ObjectID `json:"_id"`
	Name string        `json:"name"`
}

type UpdateLoadBalancerConfigRequest struct {
	Config models.LoadBalancer `json:"config"`
}

type GetProjectRoutesResponse struct {
	Routes []routeResponseModel `json:"routes"`
}

type AddRouteRequest struct {
	Path      string          `json:"path"`
	TargetURL string          `json:"target_url"`
	Method    string          `json:"method"`
	AuthMode  models.AuthMode `json:"auth_mode"`
}

type AddRouteResponse struct {
	ID bson.ObjectID `json:"_id"`
}

// hleper structs:
type routeResponseModel struct {
	ID        bson.ObjectID   `json:"_id,omitempty"`
	Path      string          `json:"path"`
	Method    string          `json:"method"`
	TargetURL string          `json:"target_url"`
	AuthMode  models.AuthMode `json:"auth_mode"`
}
