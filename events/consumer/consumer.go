package consumer

import (
	"task/consumer_service/config"
	"task/consumer_service/pkg/pubsub"

	"github.com/saidamir98/udevs_pkg/logger"
)

// Consumer ...
type Consumer struct {
	cfg     config.Config
	log     logger.LoggerI
	rmq     *pubsub.RMQ
	pubChan chan<- pubsub.Logger
}

// New ...
func New(cfg config.Config, log logger.LoggerI, rmq *pubsub.RMQ, pubChan chan<- pubsub.Logger) *Consumer {
	return &Consumer{
		cfg:     cfg,
		log:     log,
		rmq:     rmq,
		pubChan: pubChan,
	}
}

// RegisterConsumers ...
func (s *Consumer) RegisterConsumers() {
	s.rmq.AddConsumer(
		s.cfg.ExchangeName+".consumer", // consumerName
		s.cfg.ExchangeName,             // exchangeName
		s.cfg.ExchangeName+".consumer", // queueName
		s.cfg.ExchangeName+".consumer", // routingKey
		s.consumerListener,
	)
}
