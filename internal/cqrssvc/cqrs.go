package cqrssvc

import (
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/aiocean/shopify-crawler/internal/cqrssvc/router"
	"github.com/garsue/watermillzap"
	"github.com/google/wire"
	"go.uber.org/zap"
	"os"
	"strings"
)

var DefaultCqrsWireSet = wire.NewSet(
	router.DefaultRouterWireSet,
	NewCqrs,
)

func NewCqrs(
	zapLogger *zap.Logger,
	router *message.Router,
	
) (*cqrs.Facade, error) {

	logger := watermillzap.NewLogger(zapLogger)
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	eventsPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:     strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			Marshaler:   kafka.DefaultMarshaler{},
			OTELEnabled: true,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	commandsPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:     strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			Marshaler:   kafka.DefaultMarshaler{},
			OTELEnabled: true,
		},
		logger,
	)

	if err != nil {
		return nil, err
	}

	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			return commandName
		},
		GenerateEventsTopic: func(eventName string) string {
			return eventName
		},
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return kafka.NewSubscriber(
				kafka.SubscriberConfig{
					ConsumerGroup: handlerName,
					Unmarshaler:   kafka.DefaultMarshaler{},
					Brokers:       strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
					OTELEnabled:   true,
				},
				logger,
			)
		},
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return kafka.NewSubscriber(
				kafka.SubscriberConfig{
					Unmarshaler:   kafka.DefaultMarshaler{},
					ConsumerGroup: handlerName,
					Brokers:       strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
					OTELEnabled:   true,
				},
				logger,
			)
		},
		CommandEventMarshaler: cqrs.JSONMarshaler{},
		CommandsPublisher:     commandsPublisher,
		EventsPublisher:       eventsPublisher,
		Router:                router,
		Logger:                logger,
		CommandHandlers: func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{}
		},
		EventHandlers: func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler {
			return []cqrs.EventHandler{}
		},
	})

	if err != nil {
		return nil, err
	}

	return cqrsFacade, nil
}
