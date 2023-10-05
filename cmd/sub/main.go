package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/http-server/handlers/get"
	"L0/internal/http-server/handlers/get_with_url"
	mwLogger "L0/internal/http-server/middleware"
	"L0/internal/stan/stan_sub"
	"L0/internal/storage/db"
	"os/signal"

	// order "L0/internal/strct"
	pg "L0/pkg/client/postgresql"
	"L0/pkg/logger"
	"L0/pkg/logger/sl"
	"context"
	"log/slog"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting sub", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// making cache
	cache := cache.NewCache()

	// making db connection
	ctx := context.TODO()
	pg_pool, err := pg.NewClient(ctx, 1, cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	rep := db.NewRepository(pg_pool, log)
	// _ = rep
	cache.Restore(ctx, rep)
	// log.Info("\nfirst test of cache\n", cache.Items)
	// got, err := rep.GetAll(ctx)
	// if err != nil {
	// 	log.Error("couldn't get data from db", err)
	// }
	// log.Info("\nTest of GetAll func from db", got)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/", get.New(ctx, log, cache))
	// router.Post("/", get.New(ctx, log, rep))
	router.Route("/order_uid", func(r chi.Router) {
		r.Get("/{order_uid}", get_with_url.New(ctx, log, cache))
		r.Get("/", get.New(ctx, log, cache))

	})

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// if err := srv.ListenAndServe(); err != nil {
	// 	log.Error("failed to start server")
	// }
	go stan_sub.SubscribeWithParams(ctx, log, cache, rep)

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, closing server...\n\n")
			srv.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
