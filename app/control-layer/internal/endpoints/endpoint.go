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
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetUpEndpoints(pub *config.Publisher, db *mongo.Database, reg *registry.PluginRegistry) chi.Router {
	authHandler := auth.NewHandler(auth.NewRepository(db))
	projectHandler := project.NewHandler(project.NewRepository(db), reg)
	upstreamHandler := upstream.NewHandler(upstream.NewRepository(db), pub)
	routeHandler := route.NewHandler(route.NewRepository(db), pub)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Mount("/auth", authHandler)
		r.Route("/project", func(r chi.Router) {
			r.Use(
				middleware.CookieToBearer,
				middleware.Authenticator(),
				middleware.Verifier(),
			)
			r.Mount("/", projectHandler)
			r.Route("/{projectID}", func(r chi.Router) {
				r.Mount("/upstreams", upstreamHandler)
				r.Mount("/routes", routeHandler)
			})
		})
	})
	return r
}
