package project

import (
	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreatProjectRequest struct {
	Name       string   `json:"name"`
	AccessList []Access `json:"access_list"`
}

type CreatProjectResponse struct {
	Id       bson.ObjectID `json:"project_id"`
	ApiGwKey string        `json:"api_gw_key"`
}

type UpdateMiddlewaresRequest struct {
	Middlewares []models.Middleware `json:"middlewares"`
}

type UpdateAccessListRequest struct {
	AccessList []Access `json:"access_list"`
}

type GetProjectResponse struct {
	ID          bson.ObjectID       `json:"project_id"`
	Name        string              `json:"name"`
	Middlewares []models.Middleware `json:"middlewares"`
	Permission  string              `json:"permission"`
}

type ListProjectResponse struct {
	ID   bson.ObjectID `json:"project_id"`
	Name string        `json:"name"`
}

type Access struct {
	Username   string            `json:"username"`
	Permission models.Permission `json:"permission"`
}

type GetAllProjectResponse struct {
	Projects []ListProjectResponse `json:"projects"`
}
