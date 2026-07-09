package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthMode string

const (
	AuthNone     AuthMode = "none"
	AuthRequired AuthMode = "required"
)

type Route struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	ProjectID  bson.ObjectID `bson:"project_id"`
	Path       string        `bson:"path"`
	Method     string        `bson:"method"`
	UpstreamID bson.ObjectID `bson:"service_id"`
	CreatedAt  time.Time     `bson:"created_at"`
	AuthMode   AuthMode      `bson:"auth_mode"`
}
