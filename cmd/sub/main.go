package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/http-server/handlers/get"
	mwLogger "L0/internal/http-server/middleware"
	"L0/internal/storage/db"
	pg "L0/pkg/client/postgresql"
	"L0/pkg/logger"
	"L0/pkg/logger/sl"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting sub", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// making cache
	Cache := cache.NewCache()
	// cache.items[""] = order.Data{}
	// cache.data["privet"] = "jack"
	// ma, _ := json.Marshal(cache.data)
	// fmt.Println(cache, cache.data, string(ma))

	// making db connection
	ctx := context.TODO()
	pg_pool, err := pg.NewClient(ctx, 1, cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	rep := db.NewRepository(pg_pool, log)

	fmt.Println("\nfirst test of cache", Cache.Items)
	Cache.Restore(ctx, rep)
	fmt.Println("\nsecond test of cache", Cache.Items)

	// got_data, err := rep.Get(ctx, "template")
	// if err != nil {
	// 	log.Error("failed to get from db", sl.Err(err))
	// }
	// fmt.Println("got_data=", got_data)

	// got_m_data, err := rep.GetAll(ctx)
	// if err != nil {
	// 	log.Error("failed to get from db", sl.Err(err))
	// }
	// fmt.Println("\ngot_m_data=", got_m_data)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))

	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("", func(r chi.Router) {
		r.Post("/{order_uid}", get.New(ctx, log, rep))
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	clusterID := "L0_cluster"
	clientID := "L0_sub"
	URL := stan.DefaultNatsURL
	userCreds := ""

	opts := []nats.Option{nats.Name("NATS Streaming Example Publisher")}
	// Use UserCredentials
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Error("", sl.Err(err))
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Error("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	// Simple Async Subscriber

	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	}, stan.DeliverAllAvailable())
	if err != nil {
		log.Error("Subscription to Stan wasn't successful", sl.Err(err))
	}

	<-time.After(time.Second * 1)
	sub.Unsubscribe()

}
