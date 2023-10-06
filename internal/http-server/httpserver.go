package httpserver

import (
	"L0/internal/config"
	"L0/internal/http-server/handlers/get"
	"L0/internal/http-server/handlers/get_with_url"
	mwLogger "L0/internal/http-server/middleware"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// cache or db rep entity
type OrderDataGetter interface {
	Get(ctx context.Context, order_uid string) (any, error)
}

func StartServer(ctx context.Context, cfg *config.Config, log *slog.Logger, cache, rep OrderDataGetter) {
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

	// cfg.Address = "127.0.0.1:8088"
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, closing server...\n\n")
			srv.Close()
			cleanupDone <- true
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("failed to start http server", err)
	}
	<-cleanupDone

}
