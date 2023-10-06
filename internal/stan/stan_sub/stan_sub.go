package stan_sub

import (
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

	clusterID := "L0_cluster"
	clientID := "L0_sub"
	URL := stan.DefaultNatsURL
	userCreds := ""
	channel := "L0_chan"

	opts := []nats.Option{nats.Name("NATS Streaming Example Publisher")}
	// can use UserCredentials if needed
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

	// connecting to nats
	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Error("", sl.Err(err))
	}

	//connecting to stan with nats connection
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Error("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}

	// initializing simple async subscriber and describing msgHandler
	sub, err := sc.Subscribe(channel, func(m *stan.Msg) {
		log.Info("got smth!")
		log.Info("msg.Data: ", m.Data)
		var data order.Data
		data.Scan(m.Data)
		// cache.Insert(ctx, order.Data{OrderUID: "stan_inserting_into_cache"})
		// rep.Insert(ctx, order.Data{OrderUID: "stan_inserting_into_db"})
		cache.Insert(ctx, data)
		rep.Insert(ctx, data)
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
