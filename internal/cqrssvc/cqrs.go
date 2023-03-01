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
	NewCqrsSvc,
)

type CqrsSvc struct {
	cqrsFacade      *cqrs.Facade
	zapLogger       *zap.Logger
	router          *message.Router
	eventHandlers   []cqrs.EventHandler
	commandHandlers []cqrs.CommandHandler
}

func NewCqrsSvc(
	zapLogger *zap.Logger,
	router *message.Router,
) (*CqrsSvc, func(), error) {
	return &CqrsSvc{
		zapLogger: zapLogger,
		router:    router,
	}, func() {}, nil
}

func (c *CqrsSvc) AddCommandHandler(handler cqrs.CommandHandler) {
	c.commandHandlers = append(c.commandHandlers, handler)
}

func (c *CqrsSvc) AddEventHandler(handler cqrs.EventHandler) {
	c.eventHandlers = append(c.eventHandlers, handler)
}

func (c *CqrsSvc) listCommandHandlers(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler {
	return c.commandHandlers
}

func (c *CqrsSvc) listEventHandlers(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler {
	return c.eventHandlers
}

func (c *CqrsSvc) GetFacade() (*cqrs.Facade, error) {
	if c.cqrsFacade == nil {
		var err error
		c.cqrsFacade, err = c.newFacade()
		if err != nil {
			return nil, err
		}
	}

	return c.cqrsFacade, nil
}

func (c *CqrsSvc) newFacade() (*cqrs.Facade, error) {
	logger := watermillzap.NewLogger(c.zapLogger)
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
		Router:                c.router,
		Logger:                logger,
		CommandHandlers:       c.listCommandHandlers,
		EventHandlers:         c.listEventHandlers,
	})

	if err != nil {
		return nil, err
	}

	return cqrsFacade, nil
}
