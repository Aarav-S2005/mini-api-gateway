package auth

import (
	"context"
	"errors"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) users() *mongo.Collection {
	return r.db.Collection("users")
}

func (r *Repository) findUserByUsername(ctx context.Context, username string) (models.Users, error) {
	userCollection := r.users()
	var user models.Users
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Users{}, app_error.NotFound("username not found", err)
		}
		return models.Users{}, app_error.InternalServer(err)
	}
	return user, nil
}

func (r *Repository) addUser(ctx context.Context, username, password string) error {
	userCollection := r.users()
	_, err := userCollection.InsertOne(ctx, bson.M{"username": username, "password": password})
	if err != nil {
		return app_error.InternalServer(err)
	}
	return nil
}
