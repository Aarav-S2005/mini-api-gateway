package route

import (
	"context"
	"errors"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) createRoute(ctx context.Context, reqBody CreateOrUpdateRouteRequestDTO, userID, projectID bson.ObjectID) (CreateRouteResponseDTO, error) {
	err := validateDTO(reqBody)
	if err != nil {
		return CreateRouteResponseDTO{}, err
	}
	routeID, err := s.repo.createRoute(ctx, reqBody, userID, projectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return CreateRouteResponseDTO{}, app_error.NotFound("Not found", err)
		}
		if errors.Is(err, ErrRouteAlreadyExists) {
			return CreateRouteResponseDTO{}, app_error.Conflict("Route already exists", err)
		}
		return CreateRouteResponseDTO{}, err
	}
	return CreateRouteResponseDTO{RouteID: routeID}, nil
}

func (s *Service) getRouteByID(ctx context.Context, userID, projectID, routeID bson.ObjectID) (GetRouteResponseDTO, error) {
	route, err := s.repo.getRouteByID(ctx, userID, projectID, routeID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return GetRouteResponseDTO{}, app_error.NotFound("Not found", err)
		}
		return GetRouteResponseDTO{}, err
	}
	return GetRouteResponseDTO{
		ID:           route.ID,
		Path:         route.Path,
		PathType:     route.PathType,
		Method:       route.Method,
		AuthMode:     route.AuthMode,
		UpstreamName: route.UpstreamName,
		Enabled:      route.Enabled,
	}, nil
}

func (s *Service) getAllRoute(ctx context.Context, userID, projectID bson.ObjectID) (GetAllRoutesResponseDTO, error) {
	routes, err := s.repo.getAllRoutesByProjectID(ctx, userID, projectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return GetAllRoutesResponseDTO{}, app_error.NotFound("Not found", err)
		}
		return GetAllRoutesResponseDTO{}, err
	}
	routesDTO := make([]GetRouteResponseDTO, 0, len(routes))
	for _, route := range routes {
		routesDTO = append(routesDTO, GetRouteResponseDTO{
			ID:           route.ID,
			Path:         route.Path,
			PathType:     route.PathType,
			Method:       route.Method,
			AuthMode:     route.AuthMode,
			UpstreamName: route.UpstreamName,
			Enabled:      route.Enabled,
		})
	}
	return GetAllRoutesResponseDTO{Routes: routesDTO}, nil
}

func (s *Service) updateRoute(ctx context.Context, reqBody CreateOrUpdateRouteRequestDTO, userID, projectID, routeID bson.ObjectID) error {
	err := validateDTO(reqBody)
	if err != nil {
		return err
	}
	err = s.repo.updateRoute(ctx, reqBody, userID, projectID, routeID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app_error.NotFound("Not found", err)
		}
		return err
	}
	return nil
}

func (s *Service) deleteRoute(ctx context.Context, userID, projectID, routeID bson.ObjectID) error {
	err := s.repo.deleteRoute(ctx, userID, projectID, routeID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app_error.NotFound("Not found", err)
		}
		return err
	}
	return nil
}

// HELPER
func validateDTO(req CreateOrUpdateRouteRequestDTO) error {
	switch {
	case req.UpstreamName == "":
		return app_error.BadRequest("upstream_name is required", nil)
	case req.Path == "":
		return app_error.BadRequest("path is required", nil)
	case req.PathType == "":
		return app_error.BadRequest("path_type is required", nil)
	case req.Method == "":
		return app_error.BadRequest("method is required", nil)
	case req.AuthMode == "":
		return app_error.BadRequest("auth_mode is required", nil)
	}
	return nil
}
