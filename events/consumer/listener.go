package consumer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task/consumer_service/config"
	"task/consumer_service/modules"
	"task/consumer_service/pkg/helper"
	"task/consumer_service/pkg/pubsub"
	"time"

	"github.com/saidamir98/udevs_pkg/logger"
	"github.com/streadway/amqp"
)

// createProductListener ...
func (s *Consumer) consumerListener(delivery amqp.Delivery) {
	var (
		entity  modules.Consumer
		errType string
		err     error
		code    int
		resp    interface{}
	)

	count, _ := delivery.Headers["count"].(int32)
	if count > 5 { // try 5 times to consume message and if not success, then message will be dropped
		s.log.Error("message dropped", logger.String("message", string(delivery.Body)))
		return
	}

	s.pubChan <- pubsub.Logger{
		Exchange:   s.cfg.ExchangeName,
		RoutingKey: config.DebugConsumer,
		Data:       string(delivery.Body),
		Message:    "message received, time: " + time.Now().Format("2006-01-02 15:04:05"),
	}

	defer func() {
		if err != nil {
			if errType == config.ErrorConsumer {
				s.pubChan <- pubsub.Logger{
					Exchange:   s.cfg.ExchangeName,
					RoutingKey: config.ErrorConsumer,
					Data:       string(delivery.Body),
					Message:    err.Error(),
				}
			} else {
				s.pubChan <- pubsub.Logger{
					Exchange:   s.cfg.ExchangeName,
					RoutingKey: ".consumer",
					Data:       string(delivery.Body),
					Message:    err.Error(),
					Count:      int(count),
				}
			}
		}
	}()

	err = json.Unmarshal(delivery.Body, &entity)
	if err != nil {
		errType = config.ErrorConsumer
		err = fmt.Errorf("amqp message unmarshal error: %w", err)
		return
	}

	code, resp, err = helper.MakeRequest(http.MethodGet, s.cfg.RestServiceURL+"/phone/"+entity.RecordID, nil)
	if err != nil {
		err = fmt.Errorf("server doesn't ready to answer: %w", err)
		return
	}

	if code == http.StatusNotFound { // if phone not found
		errType = config.ErrorConsumer
		err = fmt.Errorf("record with id %s not found", entity.RecordID)
		return
	}

	var phone modules.Response
	err = helper.JsonToJson(&phone, resp)
	if err != nil || code != http.StatusOK { // if code is not 200 or error is not nil then the initial message is sent for reprocessing
		err = fmt.Errorf("server error or unmarshallling error, err: %w", err)
		return
	}

	js, err := json.Marshal(phone.Data)
	if err != nil {
		err = fmt.Errorf("error while marshalling: %w", err)
		return
	}

	s.pubChan <- pubsub.Logger{
		Exchange:   s.cfg.ExchangeName,
		RoutingKey: config.InfoConsumer,
		Data:       string(js),
		Message:    "phone found",
	}
}
