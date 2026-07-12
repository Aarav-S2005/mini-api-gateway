package auth

import (
	"context"
	"errors"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/app_error"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) loginUser(ctx context.Context, req RequestDTO) (string, error) {
	username, password := req.Username, req.Password
	user, err := s.repo.findUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", app_error.NotFound("Not found", err)
		}
		return "", app_error.InternalServer(err)
	}
	ok := VerifyPassword(user.Password, password)
	if !ok {
		return "", app_error.Unauthorized("Wrong credentials", err)
	}
	token, err := middleware.SignJwt(user.ID)
	if err != nil {
		return "", app_error.InternalServer(err)
	}
	return token, nil
}

func (s *Service) signUpUser(ctx context.Context, req RequestDTO) (string, error) {
	username, password := req.Username, req.Password
	_, err := s.repo.findUserByUsername(ctx, username)
	switch {
	case err == nil:
		return "", app_error.Conflict("Username already exists", nil)
	case !errors.Is(err, mongo.ErrNoDocuments):
		return "", app_error.InternalServer(err)
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", app_error.InternalServer(err)
	}
	id, err := s.repo.addUser(ctx, username, hashedPassword)
	if err != nil {
		return "", app_error.InternalServer(err)
	}
	token, err := middleware.SignJwt(id)
	if err != nil {
		return "", app_error.InternalServer(err)
	}
	return token, nil
}
