package project

import (
	"context"
	"errors"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) getCollection() *mongo.Collection {
	return repo.db.Collection("projects")
}

func (repo *Repository) createNewProject(ctx context.Context, req CreatProjectRequest, ownerID bson.ObjectID, gatewayApiKey string) (bson.ObjectID, error) {
	projects := repo.getCollection()
	project := models.Project{
		Name:               req.Name,
		OwnerId:            ownerID,
		GatewayApiKey:      gatewayApiKey,
		CreatedAt:          time.Now().UTC(),
		AccessList:         req.AccessList,
		Middlewares:        []models.Middleware{},
		LoadBalancerConfig: models.LoadBalancer{},
	}
	inserted, err := projects.InsertOne(ctx, project)
	if err != nil {
		return bson.ObjectID{}, err
	}
	id := inserted.InsertedID.(bson.ObjectID)
	return id, nil
}

func (repo *Repository) getProject(ctx context.Context, projectID, userID bson.ObjectID) (*models.Project, error) {
	projects := repo.getCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list.user_id": userID,
			},
		},
	}
	var project models.Project
	err := projects.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (repo *Repository) deleteProject(ctx context.Context, projectID, ownerID bson.ObjectID) error {
	projects := repo.getCollection()
	filter := bson.M{
		"_id":      projectID,
		"owner_id": ownerID,
	}
	_, err := projects.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) getAllProjects(ctx context.Context, userID bson.ObjectID) ([]models.Project, error) {
	projectCollection := repo.getCollection()
	filter := bson.M{
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list.user_id": userID,
			},
		},
	}
	var projects []models.Project
	cursor, err := projectCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (repo *Repository) updateAccessList(ctx context.Context, projectID, ownerID bson.ObjectID, req UpdateAccessListRequest) error {
	projectCollection := repo.getCollection()
	filter := bson.M{
		"_id":      projectID,
		"owner_id": ownerID,
	}
	update := bson.M{
		"$set": bson.M{
			"access_list": req.AccessList,
		},
	}
	_, err := projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) updateMiddlewares(ctx context.Context, projectID, userID bson.ObjectID, req UpdateMiddlewaresRequest) error {
	projectCollection := repo.getCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list.user_id":    userID,
				"access_list.permission": models.PermissionEditing,
			},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"middlewares": req.Middlewares,
		},
	}
	_, err := projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) updateOneMiddleware(ctx context.Context, projectID, userID bson.ObjectID, middleware models.Middleware) error {
	projectCollection := repo.getCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list.user_id":    userID,
				"access_list.permission": models.PermissionEditing,
			},
		},
		"middlewares.name": middleware.Name,
	}
	update := bson.M{
		"$set": bson.M{
			"middlewares.$": middleware,
		},
	}
	updated, err := projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updated.MatchedCount == 0 {
		filter = bson.M{
			"_id": projectID,
			"$or": bson.A{
				bson.M{
					"owner_id": userID,
				},
				bson.M{
					"access_list.user_id":    userID,
					"access_list.permission": models.PermissionEditing,
				},
			},
		}
		update = bson.M{
			"$push": bson.M{
				"middlewares": middleware,
			},
		}
		updated, err = projectCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		if updated.MatchedCount == 0 {
			return errors.New("no project or forbidden user")
		}
	}
	return nil
}

func (repo *Repository) deleteMiddleware(ctx context.Context, projectID, userID bson.ObjectID, name string) error {
	projectCollection := repo.getCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list.user_id":    userID,
				"access_list.permission": models.PermissionEditing,
			},
		},
	}
	update := bson.M{
		"$pull": bson.M{
			"middlewares": bson.M{
				"name": name,
			},
		},
	}
	_, err := projectCollection.UpdateOne(ctx, filter, update)
	return err
}

func (repo *Repository) deleteLoadBalancerConfig(ctx context.Context, projectID, ownerID bson.ObjectID) error {
	projectCollection := repo.getCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": ownerID,
			},
			bson.M{
				"access_list.user_id":    ownerID,
				"access_list.permission": models.PermissionEditing,
			},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"load_balancer_config": bson.Null{},
		},
	}
	_, err := projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) updateLoadBalancerConfig(ctx context.Context, projectID, ownerID bson.ObjectID, lb models.LoadBalancer) error {
	projectCollection := repo.getCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": ownerID,
			},
			bson.M{
				"access_list.user_id":    ownerID,
				"access_list.permission": models.PermissionEditing,
			},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"load_balancer_config": lb,
		},
	}
	_, err := projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
