package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Access struct {
	UserId bson.ObjectID `bson:"user_id"`
	Role   string        `bson:"role"`
}

type Project struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	Name          string        `bson:"name"`
	OwnerId       bson.ObjectID `bson:"owner_id"`
	GatewayApiKey string        `bson:"gateway_api_key"`
	CreatedAt     time.Time     `bson:"created_at"`
	AccessList    []Access      `bson:"access_list"`
}
