package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aarav-S2005/mini-api-gateway/app/db"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/config"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/health"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/lb"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/store"
	"github.com/Aarav-S2005/mini-api-gateway/app/gateway/internal/sync"
	"github.com/Aarav-S2005/mini-api-gateway/app/plugin-manager/registry"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	PluginRegistry := registry.NewRegistry()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print("ENV LOADED")

	mongocClient, err := db.NewMongo(true, cfg.MongoUri, cfg.MongoDb)
	if err != nil {
		log.Fatal(err)
		return
	}
	mongoDB := mongocClient.Database(cfg.MongoDb)
	log.Print("MONGODB UP")
	defer mongocClient.Disconnect(ctx)

	lbManager := lb.NewLBManager()
	log.Print("Load-balancer UP")

	transport := http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}
	log.Print("Transport UP")
	defer transport.CloseIdleConnections()

	initialSnapshot, err := store.LoadSnapshot(mongoDB, &transport, lbManager, PluginRegistry)
	if err != nil {
		log.Fatal(err)
		return
	}
	snapshotRegistry := store.NewRegistry(initialSnapshot)
	log.Print("SNAPSHOT REGISTRY CREATED, INITIAL SNAPSHOT BUILT")

	sl := sync.NewSynchronizeLoader(snapshotRegistry, mongoDB, &transport, lbManager, PluginRegistry)
	log.Print("SYNC LOADER CREATED")

	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUri,
		DB:       0,
		Password: "",
	})
	if err := redis.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
		return
	}
	log.Print("REDIS UP")
	sub := config.NewSubscriber(redis)
	sub.Subscribe(ctx, sl.Handler)
	defer redis.Close()

	checker := health.NewChecker(&transport, snapshotRegistry, lbManager)
	checker.Run(ctx)
}
