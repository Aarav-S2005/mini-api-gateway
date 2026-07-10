package route

import (
	"context"
	"errors"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrRouteAlreadyExists = errors.New("route already exists")

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) getRouteCollection() *mongo.Collection {
	return repo.db.Collection("route")
}

func (repo *Repository) projectAccessFilter(giveOnlyEditingPermission bool, userID, projectID bson.ObjectID) bson.M {
	permissionFilter := bson.M{
		"access_list": bson.M{
			"$elemMatch": bson.M{
				"user_id": userID,
			},
		},
	}

	if giveOnlyEditingPermission {
		permissionFilter = bson.M{
			"access_list": bson.M{
				"$elemMatch": bson.M{
					"user_id":    userID,
					"permission": models.PermissionEditing,
				},
			},
		}
	}

	return bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			permissionFilter,
		},
	}
}

func (repo *Repository) getRouteByID(ctx context.Context, userID, projectID, routeID bson.ObjectID) (*models.Route, error) {
	projectFilter := repo.projectAccessFilter(false, userID, projectID)
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return nil, err
	}
	routeFilter := bson.M{
		"_id":        routeID,
		"project_id": projectID,
	}
	var route models.Route
	err = repo.getRouteCollection().FindOne(ctx, routeFilter).Decode(&route)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

func (repo *Repository) getAllRoutesByProjectID(ctx context.Context, userID, projectID bson.ObjectID) ([]models.Route, error) {
	projectFilter := repo.projectAccessFilter(false, userID, projectID)
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return nil, err
	}
	cursor, err := repo.getRouteCollection().Find(ctx, bson.M{
		"project_id": projectID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var routes []models.Route
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}
	return routes, nil
}

func (repo *Repository) createRoute(ctx context.Context, reqBody CreateOrUpdateRouteRequestDTO, userID, projectID bson.ObjectID) (bson.ObjectID, error) {
	projectFilter := repo.projectAccessFilter(true, userID, projectID)
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return bson.ObjectID{}, err
	}
	routeFilter := bson.M{
		"project_id":    projectID,
		"path":          reqBody.Path,
		"path_type":     reqBody.PathType,
		"method":        reqBody.Method,
		"upstream_name": reqBody.UpstreamName,
	}
	err = repo.getRouteCollection().FindOne(ctx, routeFilter).Err()
	if err == nil {
		return bson.ObjectID{}, ErrRouteAlreadyExists
	}
	if !errors.Is(mongo.ErrNoDocuments, err) {
		return bson.ObjectID{}, err
	}
	now := time.Now()
	route := models.Route{
		ProjectID:    projectID,
		Path:         reqBody.Path,
		PathType:     reqBody.PathType,
		Method:       reqBody.Method,
		UpstreamName: reqBody.UpstreamName,
		AuthMode:     reqBody.AuthMode,
		Enabled:      reqBody.Enabled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	inserted, err := repo.getRouteCollection().InsertOne(ctx, route)
	if err != nil {
		return bson.ObjectID{}, err
	}

	return inserted.InsertedID.(bson.ObjectID), nil
}

func (repo *Repository) updateRoute(ctx context.Context, reqBody CreateOrUpdateRouteRequestDTO, userID, projectID, routeID bson.ObjectID) error {
	projectFilter := repo.projectAccessFilter(true, userID, projectID)
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return err
	}
	routeFilter := bson.M{
		"_id":        routeID,
		"project_id": projectID,
	}
	update := bson.M{
		"$set": bson.M{
			"updated_at":    time.Now(),
			"path":          reqBody.Path,
			"path_type":     reqBody.PathType,
			"method":        reqBody.Method,
			"upstream_name": reqBody.UpstreamName,
			"enabled":       reqBody.Enabled,
			"auth_mode":     reqBody.AuthMode,
		},
	}

	updated, err := repo.getRouteCollection().UpdateOne(ctx, routeFilter, update)
	if err != nil {
		return err
	}
	if updated.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *Repository) deleteRoute(ctx context.Context, userID, projectID, routeID bson.ObjectID) error {
	projectFilter := bson.M{
		"_id":      projectID,
		"owner_id": userID,
	}
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return err
	}
	routeFilter := bson.M{
		"_id":        routeID,
		"project_id": projectID,
	}
	deleted, err := repo.getRouteCollection().DeleteOne(ctx, routeFilter)
	if err != nil {
		return err
	}
	if deleted.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
