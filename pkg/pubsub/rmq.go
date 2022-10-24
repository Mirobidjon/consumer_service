package pubsub

import (
	"context"
	"fmt"
	"sync"

	"github.com/saidamir98/udevs_pkg/logger"

	"github.com/streadway/amqp"
)

// RMQ ...
type RMQ struct {
	log            logger.LoggerI
	amqpURI        string
	conn           *amqp.Connection
	connErr        chan *amqp.Error
	consumers      map[string]*Consumer
	consumerErrors chan error
	publishers     map[string]*Publisher
}

// NewRMQ ...
func NewRMQ(amqpURI string, log logger.LoggerI) (*RMQ, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	rmq := &RMQ{
		log:            log,
		amqpURI:        amqpURI,
		conn:           conn,
		connErr:        make(chan *amqp.Error),
		consumers:      make(map[string]*Consumer),
		consumerErrors: make(chan error, 10000), // must be buffered size
		publishers:     make(map[string]*Publisher),
	}

	rmq.conn.NotifyClose(rmq.connErr)

	go rmq.receiveConnErr()

	go rmq.receiveConsumerError()

	return rmq, nil
}

func (rmq *RMQ) receiveConnErr() {
	for err := range rmq.connErr {
		// rmq.log.Error("connection error", logger.Error(err))
		// rmq.reconnect()
		rmq.log.Panic("connection error", logger.Error(err))
	}
}

// RunConsumers ...
func (rmq *RMQ) RunConsumers(ctx context.Context) {
	var wg sync.WaitGroup

	for _, consumer := range rmq.consumers {
		wg.Add(1)
		go func(wg *sync.WaitGroup, c *Consumer) {
			defer wg.Done()
			c.Start(ctx)
		}(&wg, consumer)
		fmt.Println("Key:", consumer.queueName, "=>", "consumer:", consumer)
	}

	wg.Wait()
}

func (rmq *RMQ) receiveConsumerError() {
	for err := range rmq.consumerErrors {
		rmq.log.Error("consumer error", logger.Error(err))
	}
}
