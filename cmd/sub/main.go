package main

import (
	"L0/internal/config"
	db "L0/internal/storage/postgres"
	"L0/pkg/logger/handlers/slogpretty"
	"L0/pkg/logger/sl"
	data "L0/pkg/strct"
	"fmt"
	"os"
	"sync"
	"time"

	"log/slog"

	// json "encoding/json"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting sub", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	cache := NewCache()
	cache.data[""] = ""
	cache.data["privet"] = "poka"
	// ma, _ := json.Marshal(cache.data)
	// fmt.Println(cache, cache.data, string(ma))

	storage, err := db.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage
	var data data.Data
	data.CustomerID = "not_a_real_id"
	fmt.Println(data.Value())
	fmt.Println("\n\n\n")
	// data := map[string]string{"order_uid": "b563feb7b2b84b6test", "track_number": "WBILMTESTTRACK", "entry": "WBIL"}
	// data1, err := json.Marshal(data)
	// if err != nil {
	// 	log.Error(fmt.Sprintf("failed to marshall data"), sl.Err(err))
	// }
	in, err := storage.Save("loh", data)
	if err != nil {
		log.Error(fmt.Sprintf("failed to save to db, %d", in), sl.Err(err))
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

	<-time.After(time.Second * 10)
	sub.Unsubscribe()

}

// Cache представляет кэш для хранения данных заказов в оперативной памяти.
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewCache созадёт новый экземпляр кэша и возвращает указатель на него
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
