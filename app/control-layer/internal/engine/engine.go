package engine

import (
	"log"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/db"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints/auth"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints/project"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/go-chi/chi/v5"
)

func Run() {

	// ENV
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}

	// DB
	mongodb, err := db.NewMongo(true, cfg.MongoUri, cfg.MongoDb)
	if err != nil {
		log.Fatal(err)
		return
	}

	// REDIS PUB
	rdb, err := config.NewRedis(cfg.RedisUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	pub := config.NewPublisher(rdb)

	// JWT
	middleware.InitAuth(cfg.JwtSecret)

	// CHI ROUTER
	r := chi.NewRouter()

	// AUTH HANDLER
	authHandler := auth.NewHandler(auth.NewRepository(mongodb))
	r.Mount("/auth", authHandler)

	// PROJECT HANDLER
	projectHandler := project.NewHandler(project.NewRepository(mongodb), *registry.NewRegistry(), pub)
	r.Mount("/project", projectHandler)
}
