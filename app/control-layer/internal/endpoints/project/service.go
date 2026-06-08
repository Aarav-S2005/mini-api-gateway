package project

import (
	"context"
	"errors"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) createProject(ctx context.Context, req CreatProjectRequest, userID bson.ObjectID) (CreatProjectResponse, error) {
	gwApiKey, err := GenerateGatewayAPIKey()
	if err != nil {
		return CreatProjectResponse{}, app_error.InternalServer(err)
	}
	projectID, err := s.repo.createNewProject(ctx, req, userID, gwApiKey)
	if err != nil {
		return CreatProjectResponse{}, app_error.InternalServer(err)
	}
	return CreatProjectResponse{
		Id:       projectID,
		ApiGwKey: gwApiKey,
	}, nil
}

func (s *Service) GetProject(ctx context.Context, projectID, userID bson.ObjectID) (GetProjectResponse, error) {
	project, err := s.repo.getProject(ctx, projectID, userID)
	if err != nil {
		return GetProjectResponse{}, app_error.InternalServer(err)
	}
	var permission string
	if project.OwnerId == userID {
		permission = "admin"
	} else {
		for _, i := range project.AccessList {
			if userID == i.UserID {
				permission = string(i.Permission)
				break
			}
		}
	}
	return GetProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		CreatedAt:   project.CreatedAt,
		Middlewares: project.Middlewares,
		Permission:  permission,
	}, nil
}

func (s *Service) DeleteProject(ctx context.Context, projectID, userID bson.ObjectID) error {
	err := s.repo.deleteProject(ctx, projectID, userID)
	if err != nil {
		return app_error.InternalServer(err)
	}
	return nil
}

func (s *Service) GetAllProjects(ctx context.Context, userID bson.ObjectID) ([]GetProjectResponse, error) {
	projects, err := s.repo.getAllProjects(ctx, userID)
	if err != nil {
		return nil, app_error.InternalServer(err)
	}
	var res []GetProjectResponse
	for _, project := range projects {
		res = append(res, GetProjectResponse{
			ID:   project.ID,
			Name: project.Name,
		})
	}
	return res, nil
}

func (s *Service) UpdateAccessList(ctx context.Context, projectID, userID bson.ObjectID, req UpdateAccessListRequest) error {
	err := s.repo.updateAccessList(ctx, projectID, userID, req)
	if err != nil {
		return app_error.InternalServer(err)
	}
	return nil
}

func (s *Service) UpdateMiddlewares(ctx context.Context, projectID, userID bson.ObjectID, req UpdateMiddlewaresRequest, reg registry.PluginRegistry) error {
	for _, mw := range req.Middlewares {
		middleware, exists := reg.Get(mw.Name)
		if !exists {
			return app_error.BadRequest("invalid middleware name", nil)
		}
		err := middleware.Validate(mw.Config)
		if err != nil {
			return app_error.BadRequest("invalid middleware config", err)
		}
		err = s.repo.updateOneMiddleware(ctx, projectID, userID, mw)
		if err != nil {
			return app_error.InternalServer(err)
		}
	}
	return nil
}

func (s *Service) deleteMiddleware(ctx context.Context, projectID, userID bson.ObjectID, name string) error {
	err := s.repo.deleteMiddleware(ctx, projectID, userID, name)
	if err != nil {
		return app_error.InternalServer(err)
	}
	return nil
}

func (s *Service) deleteLoadBalancerConfig(ctx context.Context, projectID, userID bson.ObjectID) error {
	err := s.repo.deleteLoadBalancerConfig(ctx, projectID, userID)
	if err != nil {
		return app_error.InternalServer(err)
	}
	return nil
}

func (s *Service) updateLoadBalancerConfig(ctx context.Context, projectID, userID bson.ObjectID, req UpdateLoadBalancerConfigRequest) error {
	cfg := req.Config
	if cfg.RetryAttemptsAllowed < 1 {
		return app_error.BadRequest("invalid retry attempts allowed", errors.New("invalid retry attempts allowed: 0 retry count"))
	}
	if len(cfg.Backends) == 0 {
		return app_error.BadRequest("invalid backends", errors.New("backends len is 0"))
	}
	for _, backend := range cfg.Backends {
		err := ValidateBackendURL(backend.URL)
		if err != nil {
			return app_error.BadRequest("invalid backends", err)
		}
		if backend.Weight < 1 {
			return app_error.BadRequest("invalid backends", errors.New("invalid backends: weight is less than 1"))
		}
	}
	err := s.repo.updateLoadBalancerConfig(ctx, projectID, userID, cfg)
	if err != nil {
		return app_error.InternalServer(err)
	}
	return nil
}
