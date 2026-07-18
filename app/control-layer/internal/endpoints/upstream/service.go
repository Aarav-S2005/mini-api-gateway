package upstream

import (
	"context"
	"errors"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/lib"
	"github.com/Aarav-S2005/mini-api-gateway/app/db/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) createUpstream(ctx context.Context, userID, projectID bson.ObjectID, reqBody CreateOrUpdateUpstreamRequestDTO) (bson.ObjectID, error) {
	err := validateDTO(reqBody)
	if err != nil {
		return bson.ObjectID{}, err
	}
	upstreamID, err := s.repo.createUpstream(ctx, reqBody, userID, projectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return bson.ObjectID{}, app_error.NotFound("project or upstream not found", err)
		}
		if errors.Is(err, ErrUpstreamAlreadyExists) {
			return bson.ObjectID{}, app_error.Conflict("upstream already exists", err)
		}
		return bson.ObjectID{}, err
	}
	return upstreamID, nil
}

func (s *Service) getAllUpstream(ctx context.Context, userID, projectID bson.ObjectID) (GetAllUpstreamResponseDTO, error) {
	upstreams, err := s.repo.getAllUpstreamsByProjectID(ctx, userID, projectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return GetAllUpstreamResponseDTO{}, app_error.NotFound("project or upstream not found", err)
		}
		return GetAllUpstreamResponseDTO{}, err
	}
	upstreamsDTO := make([]GetUpstreamResponseDTO, 0, len(upstreams))
	for _, upstream := range upstreams {
		upstreamsDTO = append(upstreamsDTO, GetUpstreamResponseDTO{
			ID:                    upstream.ID,
			Name:                  upstream.Name,
			LoadBalancingStrategy: upstream.LoadBalancingStrategy,
			Backends:              upstream.Backends,
		})
	}
	return GetAllUpstreamResponseDTO{
		Upstreams: upstreamsDTO,
	}, nil
}

func (s *Service) getUpstreamByID(ctx context.Context, userID, projectID, upstreamID bson.ObjectID) (GetUpstreamResponseDTO, error) {
	upstream, err := s.repo.getUpstreamById(ctx, userID, projectID, upstreamID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return GetUpstreamResponseDTO{}, app_error.NotFound("project or upstream not found", err)
		}
		return GetUpstreamResponseDTO{}, err
	}
	return GetUpstreamResponseDTO{
		ID:                    upstream.ID,
		Name:                  upstream.Name,
		LoadBalancingStrategy: upstream.LoadBalancingStrategy,
		Backends:              upstream.Backends,
	}, nil
}

func (s *Service) updateUpstream(ctx context.Context, userID, projectID, upstreamID bson.ObjectID, reqBody CreateOrUpdateUpstreamRequestDTO) error {
	err := validateDTO(reqBody)
	if err != nil {
		return err
	}
	err = s.repo.updateUpstream(ctx, reqBody, userID, projectID, upstreamID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app_error.NotFound("project or upstream not found", err)
		}
		return err
	}
	return nil
}

func (s *Service) deleteUpstream(ctx context.Context, userID, projectID, upstreamID bson.ObjectID) error {
	err := s.repo.deleteUpstream(ctx, userID, projectID, upstreamID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return app_error.NotFound("project or upstream not found", err)
		}
		return err
	}
	return nil
}

// HELPER
func validateDTO(reqBody CreateOrUpdateUpstreamRequestDTO) error {
	if reqBody.Name == "" || len(reqBody.Backends) == 0 || reqBody.LoadBalancingStrategy == "" {
		return app_error.BadRequest("invalid request body", errors.New("name, load balancing strategy and backend are required"))
	}
	for _, backend := range reqBody.Backends {
		if backend.URL == "" {
			return app_error.BadRequest("invalid data in json", errors.New("backend URL cannot be empty"))
		}

		if err := lib.ValidateBackendURL(backend.URL); err != nil {
			return app_error.BadRequest("invalid data in json", err)
		}

		if isValidLoadBalancingStrategy(reqBody.LoadBalancingStrategy) || reqBody.LoadBalancingStrategy == models.WeightedRoundRobin &&
			backend.Weight == nil {
			return app_error.BadRequest("invalid data in json", errors.New("weight is required for weighted round robin"))
		}
	}
	return nil
}

func isValidLoadBalancingStrategy(s models.LoadBalancingStrategy) bool {
	switch s {
	case models.RoundRobinLoadBalancing,
		models.RandomLoadBalancing,
		models.IPHashLoadBalancing,
		models.WeightedRoundRobin,
		models.LeastConnections:
		return true
	default:
		return false
	}
}
