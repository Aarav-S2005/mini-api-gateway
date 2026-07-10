package route

import (
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func getUserIDAndProjectID(r *http.Request) (bson.ObjectID, bson.ObjectID, error) {
	userID, err := lib.GetUserID(r.Context())
	if err != nil {
		return bson.ObjectID{}, bson.ObjectID{}, app_error.Unauthorized("invalid userID", err)
	}
	projectID, err := lib.GetIdFromEndpoint(r, "projectID")
	if err != nil {
		return bson.ObjectID{}, bson.ObjectID{}, app_error.BadRequest("no projectID present", err)
	}
	return userID, projectID, nil
}

func getUserIDProjectIDAndRouteID(r *http.Request) (bson.ObjectID, bson.ObjectID, bson.ObjectID, error) {
	userID, projectID, err := getUserIDAndProjectID(r)
	if err != nil {
		return bson.ObjectID{}, bson.ObjectID{}, bson.ObjectID{}, err
	}
	routeID, err := lib.GetIdFromEndpoint(r, "routeID")
	if err != nil {
		return bson.ObjectID{}, bson.ObjectID{}, bson.ObjectID{}, app_error.BadRequest("no routeID present", err)
	}
	return userID, projectID, routeID, nil
}
