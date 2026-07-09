package route

import (
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	pub     *config.Publisher
}

func NewHandler(repo *Repository, pub *config.Publisher) *Handler {
	return &Handler{service: NewService(repo), pub: pub}
}

func (h *Handler) initRoutes() chi.Router {
	r := chi.NewRouter()
	return r
}
