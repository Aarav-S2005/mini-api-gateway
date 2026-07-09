package route

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repository struct {
	db *mongo.Database
}
