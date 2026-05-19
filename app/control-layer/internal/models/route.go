package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Route struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	ProjectID   bson.ObjectID `bson:"project_id"`
	Path        string        `bson:"path"`
	Method      string        `bson:"method"`
	TargetURL   string        `bson:"target_url"`
	CreatedAt   time.Time     `bson:"created_at"`
	Middlewares []Middleware  `bson:"middlewares"`
}
