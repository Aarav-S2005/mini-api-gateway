package upstream

import (
	"context"
	"errors"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrUpstreamAlreadyExists = errors.New("upstream already exists")

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) getUpstreamCollection() *mongo.Collection {
	return repo.db.Collection("upstreams")
}

func (repo *Repository) projectAccessFilter(giveOnlyEditingPermission bool, userID, projectId bson.ObjectID) bson.M {
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
		"_id": projectId,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			permissionFilter,
		},
	}
}

func (repo *Repository) getUpstreamById(ctx context.Context, userID, projectID, serviceID bson.ObjectID) (*models.Upstream, error) {
	projectFilter := repo.projectAccessFilter(false, userID, projectID)
	var project models.Project
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Decode(&project)
	if err != nil {
		return nil, err
	}
	upstreamFilter := bson.M{
		"_id":        serviceID,
		"project_id": projectID,
	}
	var upstream models.Upstream
	err = repo.getUpstreamCollection().FindOne(ctx, upstreamFilter).Decode(&upstream)
	if err != nil {
		return nil, err
	}
	return &upstream, nil
}

func (repo *Repository) getAllUpstreamsByProjectID(ctx context.Context, userID, projectID bson.ObjectID) ([]models.Upstream, error) {
	projectFilter := repo.projectAccessFilter(false, userID, projectID)
	var project models.Project
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Decode(&project)
	if err != nil {
		return nil, err
	}
	cursor, err := repo.getUpstreamCollection().Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var upstreams []models.Upstream
	if err = cursor.All(ctx, &upstreams); err != nil {
		return nil, err
	}
	return upstreams, nil
}

func (repo *Repository) createUpstream(ctx context.Context, reqBody CreateOrUpdateUpstreamRequestDTO, userID, projectID bson.ObjectID) (bson.ObjectID, error) {
	projectFilter := repo.projectAccessFilter(true, userID, projectID)
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return bson.ObjectID{}, err
	}
	upstreamFilter := bson.M{
		"project_id": projectID,
		"name":       reqBody.Name,
	}
	err = repo.getUpstreamCollection().FindOne(ctx, upstreamFilter).Err()
	if err == nil {
		return bson.ObjectID{}, ErrUpstreamAlreadyExists
	}
	if !errors.Is(mongo.ErrNoDocuments, err) {
		return bson.ObjectID{}, err
	}
	now := time.Now()
	upstream := models.Upstream{
		ProjectID:             projectID,
		Name:                  reqBody.Name,
		LoadBalancingStrategy: reqBody.LoadBalancingStrategy,
		Backends:              reqBody.Backends,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	inserted, err := repo.getUpstreamCollection().InsertOne(ctx, upstream)
	if err != nil {
		return bson.ObjectID{}, err
	}
	id := inserted.InsertedID.(bson.ObjectID)
	return id, nil
}

func (repo *Repository) updateUpstream(ctx context.Context, reqBody CreateOrUpdateUpstreamRequestDTO, userID, projectID, upstreamID bson.ObjectID) error {
	projectFilter := repo.projectAccessFilter(true, userID, projectID)
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return err
	}
	upstreamFilter := bson.M{
		"_id":        upstreamID,
		"project_id": projectID,
	}
	update := bson.M{
		"$set": bson.M{
			"updated_at":              time.Now(),
			"name":                    reqBody.Name,
			"load_balancing_strategy": reqBody.LoadBalancingStrategy,
			"backends":                reqBody.Backends,
		},
	}
	updated, err := repo.getUpstreamCollection().UpdateOne(ctx, upstreamFilter, update)
	if err != nil {
		return err
	}
	if updated.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *Repository) deleteUpstream(ctx context.Context, userID, projectID, upstreamID bson.ObjectID) error {
	projectFilter := bson.M{
		"_id":      projectID,
		"owner_id": userID,
	}
	err := repo.db.Collection("project").FindOne(ctx, projectFilter).Err()
	if err != nil {
		return err
	}
	upstreamFilter := bson.M{
		"_id":        upstreamID,
		"project_id": projectID,
	}
	deleted, err := repo.getUpstreamCollection().DeleteOne(ctx, upstreamFilter)
	if err != nil {
		return err
	}
	if deleted.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
