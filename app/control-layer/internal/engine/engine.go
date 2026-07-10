package engine

import (
	"log"
	"net/http"

	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/db"
	"github.com/Aarav-S2005/mini-api-gateway/app/control-layer/internal/endpoints"
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

	// DB
	mongodb, err := db.NewMongo(true, cfg.MongoUri, cfg.MongoDb)
	if err != nil {
		log.Fatal(err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	pub := config.NewPublisher(rdb)

	r := endpoints.SetUpEndpoints(pub, mongodb, registry.NewRegistry())

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
		return
	}

}
