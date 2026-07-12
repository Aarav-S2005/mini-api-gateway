package project

import (
	"context"
	"errors"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	db *mongo.Database
}

var ErrNoUserFound = errors.New("username not found")

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) getProjectCollection() *mongo.Collection {
	return repo.db.Collection("project")
}

func (repo *Repository) createNewProject(ctx context.Context, req CreatProjectRequest, ownerID bson.ObjectID, gatewayApiKey string) (bson.ObjectID, error) {
	projects := repo.getProjectCollection()
	accessUsernames := make([]string, len(req.AccessList))
	for i, access := range req.AccessList {
		accessUsernames[i] = access.Username
	}
	filter := bson.M{
		"username": bson.M{
			"$in": accessUsernames,
		},
	}
	cursor, err := repo.db.Collection("users").Find(ctx, filter)
	if err != nil {
		return bson.ObjectID{}, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return bson.ObjectID{}, err
	}
	usernameToID := make(map[string]bson.ObjectID, len(users))
	for _, user := range users {
		usernameToID[user.Username] = user.ID
	}
	accessList := make([]models.Access, len(req.AccessList))
	for i, access := range req.AccessList {
		userID, ok := usernameToID[access.Username]
		if !ok {
			return bson.ObjectID{}, ErrNoUserFound
		}

		accessList[i] = models.Access{
			UserID:     userID,
			Permission: access.Permission,
		}
	}

	project := models.Project{
		Name:          req.Name,
		OwnerId:       ownerID,
		GatewayApiKey: gatewayApiKey,
		CreatedAt:     time.Now().UTC(),
		AccessList:    accessList,
		Middlewares:   []models.Middleware{},
	}
	inserted, err := projects.InsertOne(ctx, project)
	if err != nil {
		return bson.ObjectID{}, err
	}
	id := inserted.InsertedID.(bson.ObjectID)
	return id, nil
}

func (repo *Repository) getProject(ctx context.Context, projectID, userID bson.ObjectID) (*models.Project, error) {
	projects := repo.getProjectCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list": bson.M{
					"$elemMatch": bson.M{
						"user_id": userID,
					},
				},
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
	projects := repo.getProjectCollection()
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
	projectCollection := repo.getProjectCollection()
	filter := bson.M{
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list": bson.M{
					"$elemMatch": bson.M{
						"user_id": userID,
					},
				},
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
	projectCollection := repo.getProjectCollection()
	filter := bson.M{
		"_id":      projectID,
		"owner_id": ownerID,
	}
	accessUsernames := make([]string, len(req.AccessList))
	for i, access := range req.AccessList {
		accessUsernames[i] = access.Username
	}
	filterAL := bson.M{
		"username": bson.M{
			"$in": accessUsernames,
		},
	}
	cursor, err := repo.db.Collection("users").Find(ctx, filterAL)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return err
	}
	usernameToID := make(map[string]bson.ObjectID, len(users))
	for _, user := range users {
		usernameToID[user.Username] = user.ID
	}
	accessList := make([]models.Access, len(req.AccessList))
	for i, access := range req.AccessList {
		userID, ok := usernameToID[access.Username]
		if !ok {
			return ErrNoUserFound
		}

		accessList[i] = models.Access{
			UserID:     userID,
			Permission: access.Permission,
		}
	}
	update := bson.M{
		"$set": bson.M{
			"access_list": accessList,
		},
	}
	_, err = projectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) updateMiddlewares(ctx context.Context, projectID, userID bson.ObjectID, req UpdateMiddlewaresRequest) error {
	projectCollection := repo.getProjectCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list": bson.M{
					"$elemMatch": bson.M{
						"user_id":    userID,
						"permission": models.PermissionEditing,
					},
				},
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
	projectCollection := repo.getProjectCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list": bson.M{
					"$elemMatch": bson.M{
						"user_id":    userID,
						"permission": models.PermissionEditing,
					},
				},
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
					"access_list": bson.M{
						"$elemMatch": bson.M{
							"user_id":    userID,
							"permission": models.PermissionEditing,
						},
					},
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
	projectCollection := repo.getProjectCollection()
	filter := bson.M{
		"_id": projectID,
		"$or": bson.A{
			bson.M{
				"owner_id": userID,
			},
			bson.M{
				"access_list": bson.M{
					"$elemMatch": bson.M{
						"user_id":    userID,
						"permission": models.PermissionEditing,
					},
				},
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
