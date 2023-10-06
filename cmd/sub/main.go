package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	httpserver "L0/internal/http-server"
	"L0/internal/stan/stan_sub"
	"L0/internal/storage/db"
	pg "L0/pkg/client/postgresql"
	"L0/pkg/logger"
	"L0/pkg/logger/sl"
	"context"
	"log/slog"
	"os"
	"sync"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting sub", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx := context.TODO()

	// making db connection and initializing db repository
	pg_pool, err := pg.NewClient(ctx, 1, cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	rep := db.NewRepository(pg_pool, log)

	// making cache
	cache := cache.NewCache()
	// restoring cache from db rep
	cache.Restore(ctx, rep)

	wg := sync.WaitGroup{}

	// starting stan subscriber entity
	wg.Add(1)
	go func() {
		stan_sub.SubscribeWithParams(ctx, log, cache, rep)
		wg.Done()
	}()

	// starting http-server to handle API requests
	wg.Add(1)
	go func() {
		httpserver.StartServer(ctx, cfg, log, cache, rep)
		wg.Done()
	}()
	wg.Wait()
	log.Info("app closed")
}
