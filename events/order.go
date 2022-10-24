package events

import (
	"context"
	"task/consumer_service/config"
	"task/consumer_service/events/consumer"
	"task/consumer_service/pkg/pubsub"

	"github.com/saidamir98/udevs_pkg/logger"
)

// PubsubServer ...
type PubsubServer struct {
	cfg config.Config
	log logger.LoggerI
	RMQ *pubsub.RMQ
}

// New ...
func New(cfg config.Config, log logger.LoggerI) (*PubsubServer, error) {
	rmq, err := pubsub.NewRMQ(cfg.RabbitMqURL, log)
	if err != nil {
		return nil, err
	}

	rmq.AddPublisher(cfg.ExchangeName)
	// rmq.AddPublisher(cfg.ExchangeName + config.ErrorConsumer)
	// rmq.AddPublisher(cfg.ExchangeName + config.InfoConsumer)
	// rmq.AddPublisher(cfg.ExchangeName + config.DebugConsumer)

	return &PubsubServer{
		cfg: cfg,
		log: log,
		RMQ: rmq,
	}, nil
}

// Run ...
func (s *PubsubServer) Run(ctx context.Context, pubChan chan<- pubsub.Logger) {
	consumerServer := consumer.New(s.cfg, s.log, s.RMQ, pubChan)
	consumerServer.RegisterConsumers()

	s.RMQ.RunConsumers(ctx)
}
