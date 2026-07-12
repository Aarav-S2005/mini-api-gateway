package engine

import (
	"log"
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/middleware"
	"github.com/Aarav-S2005/mini-api-gateway/app/db"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/redis/go-redis/v9"
)

func Run() {

	// ENV
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("env loaded")

	middleware.InitAuth(cfg.JwtSecret)

	// DB
	mongodb, err := db.NewMongo(true, cfg.MongoUri, cfg.MongoDb)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("mongodb connected")

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUri,
		DB:       0,
		Password: "",
	})
	pub := config.NewPublisher(rdb)
	log.Println("redis connected")

	r := endpoints.SetUpEndpoints(pub, mongodb, registry.NewRegistry())
	log.Println("endpoints up")

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
		return
	}

}
