package main

import (
	"L0/internal/config"
	"L0/pkg/logger/handlers/slogpretty"
	"L0/pkg/logger/sl"

	// "fmt"
	"os"
	"time"

	"log/slog"

	stan "github.com/nats-io/stan.go"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting sub", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	clusterID := "L0_cluster"
	// clusterID = stan.DefaultNatsURL
	// clusterID = "nats://127.0.0.1:4223"
	clientID := "L0_pub"

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Error("Connection to Stan wasn't successful", sl.Err(err))
	}

	// Simple Synchronous Publisher
	sc.Publish("foo", []byte("Hello 1")) // does not return until an ack has been received from NATS Streaming
	time.Sleep(time.Second * 2)
	sc.Publish("foo", []byte("Hello 2"))
	time.Sleep(time.Second * 2)
	sc.Publish("foo", []byte("Hello 3"))
	time.Sleep(time.Second * 2)
	// Simple Async Subscriber
	// sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
	// 	fmt.Printf("Received a message: %s\n", string(m.Data))
	// })
	// if err != nil {
	// 	log.Error("Subscription to Stan wasn't successful", sl.Err(err))
	// }

	// Unsubscribe
	// sub.Unsubscribe()

	// Close connection
	sc.Close()
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
		// log = slog.New(
		// 	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// )
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
