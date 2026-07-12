package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission string

const (
	PermissionEditing Permission = "editing"
	PermissionViewing Permission = "viewing"
)

type Access struct {
	UserID     bson.ObjectID `bson:"user_id" json:"user_id"`
	Permission Permission    `bson:"permission" json:"permission"`
}

type Project struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	Name          string        `bson:"name"`
	OwnerId       bson.ObjectID `bson:"owner_id"`
	GatewayApiKey string        `bson:"gateway_api_key"`
	CreatedAt     time.Time     `bson:"created_at"`
	AccessList    []Access      `bson:"access_list"`
	Middlewares   []Middleware  `bson:"middlewares"`
}
