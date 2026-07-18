package main

import (
	"log"

	"github.com/Aarav-S2005/mini-api-gateway/app/db"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/store"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
)

func main() {
	MiddlewareRegistry := registry.NewRegistry()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}

	mongoDB, err := db.NewMongo(true, cfg.MongoUri, cfg.MongoDb)
	if err != nil {
		log.Fatal(err)
		return
	}

	initialSnapshot, err := store.LoadSnapshot(mongoDB)
	if err != nil {
		log.Fatal(err)
		return
	}

	snapshotRegistry := store.NewRegistry(initialSnapshot)

}
