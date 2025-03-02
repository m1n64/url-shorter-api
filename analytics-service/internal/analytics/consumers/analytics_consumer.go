package consumers

import (
	"analytics-service/internal/analytics/entities"
	"analytics-service/internal/analytics/services"
	"analytics-service/pkg/utils"
	"github.com/go-json-experiment/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type analyticsConsumer struct {
	analyticsEventService *services.AnalyticsEventService
	logger                *zap.Logger
}

func NewAnalyticsConsumer(analyticsEventService *services.AnalyticsEventService, logger *zap.Logger) utils.Consumer {
	return &analyticsConsumer{
		analyticsEventService: analyticsEventService,
		logger:                logger,
	}
}

func (c *analyticsConsumer) Handle(msg amqp.Delivery) {
	var eventList []*entities.AnalyticsEvent

	err := json.Unmarshal(msg.Body, &eventList)
	if err != nil {
		c.logger.Error("Error unmarshalling analytics event", zap.Error(err))
		return
	}

	go func(events []*entities.AnalyticsEvent) {
		for _, event := range events {
			_, err := c.analyticsEventService.Save(event)
			if err != nil {
				c.logger.Error("Error saving analytics event", zap.Any("event", event), zap.Error(err))
				return
			}
		}

		if err != nil {
			c.logger.Error("Error saving analytics event", zap.Error(err))
		}
	}(eventList)

	err = msg.Ack(false)
	if err != nil {
		c.logger.Error("Error acking analytics event", zap.Error(err))
		return
	}
}
