package engine

import (
	"log"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	redis2 "github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config/redis"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/db"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/routes/auth"
	"github.com/go-chi/chi/v5"
)

func Run() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	mongodb, err := db.NewMongo(true, cfg.MongoUri, cfg.MongoDb)
	if err != nil {
		log.Fatal(err)
		return
	}
	rdb, err := redis2.NewRedis(cfg.RedisUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	pub := redis2.NewPublisher(rdb)
	middleware.InitAuth(cfg.JwtSecret)
	r := chi.NewRouter()
	authHandler := auth.NewHandler(auth.NewRepository(mongodb))
	r.Mount("/auth", authHandler)
}
