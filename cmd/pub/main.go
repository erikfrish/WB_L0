package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	order "L0/internal/strct"
	"L0/pkg/logger"
	"L0/pkg/logger/sl"
	"encoding/json"
	"os"

	// "fmt"

	"time"

	"log/slog"

	stan "github.com/nats-io/stan.go"
)

func main() {

	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting pub", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// making cache
	cache := cache.NewCache()
	_ = cache

	// mock_data_path := "./../../mock_orders/mockturtle.json"
	mock_data_path := "mock_orders/mockturtle.json"
	// check if file exists
	if _, err := os.Stat(mock_data_path); os.IsNotExist(err) {
		log.Error("mock data file does not exist: %s", mock_data_path)
	}
	mock_turtle, err := os.ReadFile(mock_data_path)
	if err != nil {
		log.Error("failed to read mock data file", err)
	}
	var mock_data []order.Data
	if err = json.Unmarshal(mock_turtle, &mock_data); err != nil {
		log.Error("failed to unmarshal data from file", err)
	}
	// pri, _ := json.Marshal(mock_data)
	log.Info("opening mock data file succeeded, len of data:", len(mock_data))

	clusterID := "L0_cluster"
	clientID := "L0_pub"
	channel := "L0_chan"
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Error("Connection to Stan wasn't successful", sl.Err(err))
	}

	// Simple Synchronous Publisher
	for _, v := range mock_data {
		msg_to_send, err := json.Marshal(v)
		if err != nil {
			log.Error("failed to unmarshal mock data", err)
		}
		sc.Publish(channel, msg_to_send)
		time.Sleep(time.Second * 1)
	}
	// Close connection
	sc.Close()
}
