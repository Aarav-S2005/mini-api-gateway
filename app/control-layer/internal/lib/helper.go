package lib

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetUserID(ctx context.Context) (bson.ObjectID, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return bson.ObjectID{}, err
	}
	userIDStr, ok := claims["userID"].(string)
	if !ok {
		return bson.ObjectID{}, errors.New("userID claim missing or invalid")
	}
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return userID, nil
}

func GetIdFromEndpoint(r *http.Request, key string) (bson.ObjectID, error) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, key))
	if err != nil {
		return bson.ObjectID{}, err
	}
	return id, nil
}

func GetProjectAndUserID(r *http.Request) (bson.ObjectID, bson.ObjectID, error) {
	userID, err := GetUserID(r.Context())
	if err != nil {
		return bson.ObjectID{0}, bson.ObjectID{}, err
	}
	projectID, err := GetIdFromEndpoint(r, "projectID")
	if err != nil {
		return bson.ObjectID{}, bson.ObjectID{}, err
	}
	return projectID, userID, nil
}
