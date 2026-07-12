package endpoints

import (
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints/auth"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints/project"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints/route"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints/upstream"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/go-chi/chi/v5"
	chiMid "github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetUpEndpoints(pub *config.Publisher, db *mongo.Database, reg *registry.PluginRegistry) chi.Router {
	authHandler := auth.NewHandler(auth.NewRepository(db))
	upstreamHandler := upstream.NewHandler(upstream.NewRepository(db), pub)
	routeHandler := route.NewHandler(route.NewRepository(db), pub)
	projectHandler := project.NewHandler(project.NewRepository(db), reg, upstreamHandler, routeHandler)

	r := chi.NewRouter()
	r.Use(chiMid.Logger)
	r.Use(chiMid.Recoverer)
	r.Group(func(r chi.Router) {
		r.Mount("/auth", authHandler)
		r.Route("/projects", func(r chi.Router) {
			r.Use(
				middleware.CookieToBearer,
				middleware.Verifier(),
				middleware.Authenticator(),
			)
			r.Mount("/", projectHandler)
		})
	})
	return r
}
