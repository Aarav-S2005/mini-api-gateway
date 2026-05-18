package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Users struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
}
