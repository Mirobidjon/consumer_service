package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"task/consumer_service/config"
	"task/consumer_service/events"
	"task/consumer_service/pkg/pubsub"

	"github.com/saidamir98/udevs_pkg/logger"
)

func main() {
	var loggerLevel string
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
	case config.TestMode:
		loggerLevel = logger.LevelDebug
	default:
		loggerLevel = logger.LevelInfo
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)

	pubChan := make(chan pubsub.Logger, 100) // 100 is enough for application service

	pubsubServer, err := events.New(cfg, log)
	if err != nil {
		log.Panic("error on the event server", logger.Error(err))
	}

	go func() { // run publisher channel
		err = pubsub.RunPublisherChannel(pubChan, pubsubServer.RMQ)
		if err != nil {
			log.Panic("error on the publisher channel", logger.Error(err))
		}

		close(pubChan)
	}()

	go func() {
		pubsubServer.Run(ctx, pubChan)
	}()

	log.Info("service is running")

	shutdownChan := make(chan os.Signal, 1)
	defer close(shutdownChan)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-shutdownChan

	cancel()

	log.Info("received os signal", logger.Any("signal", sig))
	// wait for publisher channel to close
	close(pubChan)

	log.Info("server shutdown successfully")
}
