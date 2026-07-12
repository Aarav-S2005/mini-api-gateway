package snapshot

import (
	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ProjectID   bson.ObjectID
	Middlewares []models.Middleware
}
