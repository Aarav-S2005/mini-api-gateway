package lib

import (
	"context"
	"errors"

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
