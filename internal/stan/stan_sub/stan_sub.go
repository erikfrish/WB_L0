package stan_sub

import (
	"L0/internal/config"
	order "L0/internal/strct"
	"L0/pkg/logger/sl"
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// cache or db rep entity
type DataSaver interface {
	Insert(ctx context.Context, data order.Data) error
}

// make subscription with msgHandler, which inserts data to cache or db rep
func SubscribeWithParams(ctx context.Context, log *slog.Logger, cache, rep DataSaver) {

	cfg := config.MustLoad()

	opts := []nats.Option{nats.Name("NATS Streaming Example Publisher")}
	// can use UserCredentials if needed
	if cfg.Stan.UserCreds != "" {
		opts = append(opts, nats.UserCredentials(cfg.Stan.UserCreds))
	}

	// connecting to nats
	nc, err := nats.Connect(cfg.Stan.URL, opts...)
	if err != nil {
		log.Error("Can't connect to nats", sl.Err(err))
	}

	//connecting to stan with nats connection
	sc, err := stan.Connect(cfg.Stan.ClusterID, cfg.Stan.ClientID, stan.NatsConn(nc))
	if err != nil {
		log.Error("Can't connect to stan: %v.\nMake sure a NATS Streaming Server is running at: %s", err, cfg.Stan.URL)
	}

	// initializing simple async subscriber and describing msgHandler
	sub, err := sc.Subscribe(cfg.Stan.Channel, func(m *stan.Msg) {
		log.Debug("got smth!")
		// log.Debug("msg.Data: ", m.Data)
		var data order.Data
		if err := data.Scan(m.Data); err != nil {
			log.Error(err.Error())
		}
		// log.Debug("data decoded: ", data)
		if err := rep.Insert(ctx, data); err != nil {
			log.Error("can't Insert into db", err.Error())
		}
		if err := cache.Insert(ctx, data); err != nil {
			log.Error("can't Insert into cache", err.Error())
		}
	})
	// }, stan.StartWithLastReceived())
	// }, stan.DeliverAllAvailable())

	if err != nil {
		log.Error("Subscription to Stan wasn't successful", sl.Err(err))
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sub.Unsubscribe()
			sc.Close()
			nc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
