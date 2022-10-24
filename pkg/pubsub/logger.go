package pubsub

import (
	"task/consumer_service/config"
	"time"

	"github.com/google/uuid"
	"github.com/saidamir98/udevs_pkg/logger"
	"github.com/streadway/amqp"
)

type Logger struct {
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
	Data       string `json:"data"`
	Message    string `json:"message"`
	Count      int    `json:"count"`
}

func RunPublisherChannel(tasks <-chan Logger, rmq *RMQ) error {
	var (
		err     error
		errChan = make(chan error)
	)

	for task := range tasks {
		select {
		case err = <-errChan:
			return err
		default:
			go func(task *Logger) {
				time.Sleep(time.Duration(task.Count*10) * time.Second) // sleep for count * 10 seconds, so increase sleep time for each retry
				err = rmq.Push(
					task.Exchange,
					task.Exchange+task.RoutingKey,
					amqp.Publishing{
						ContentType:  "application/json",
						Body:         []byte(task.Data),
						Type:         "JSON",
						MessageId:    uuid.NewString(),
						DeliveryMode: amqp.Persistent,
						Headers: amqp.Table{
							"message": task.Message,
							"count":   task.Count + 1,
						},
					},
				)
				if err != nil {
					errChan <- err
					return
				}

				if task.Count == 0 {
					err = rmq.Push(
						task.Exchange,
						task.Exchange+config.AllMessageConsumer, // push to all message consumer
						amqp.Publishing{
							ContentType:  "application/json",
							Body:         []byte(task.Data),
							Type:         "JSON",
							MessageId:    uuid.NewString(),
							DeliveryMode: amqp.Persistent,
							Headers: amqp.Table{
								"message": task.Message,
							},
						},
					)
					if err != nil {
						errChan <- err
						return
					}
				}

				rmq.log.Info("message sent to exchange", logger.String("exchange", task.Exchange+task.RoutingKey), logger.String("message", task.Message))
			}(&task)
		}
	}

	return nil
}
