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
	Id       bson.ObjectID `json:"project_id"`
	ApiGwKey string        `json:"api_gw_key"`
}

type UpdateMiddlewaresRequest struct {
	Middlewares []models.Middleware `json:"middleware"`
}

type UpdateAccessListRequest struct {
	AccessList []models.Access `json:"access_list"`
}

type GetProjectResponse struct {
	ID          bson.ObjectID       `json:"project_id"`
	Name        string              `json:"name"`
	CreatedAt   time.Time           `json:"created_at"`
	Middlewares []models.Middleware `json:"middlewares"`
	Permission  string              `json:"permission"`
}

type ListProjectResponse struct {
	ID   bson.ObjectID `json:"project_id"`
	Name string        `json:"name"`
}
