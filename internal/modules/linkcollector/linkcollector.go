package linkcollector

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/wire"
)

var DefaultLinkCollectorWireSet = wire.NewSet(
	NewLinkCollectorModule,
	DefaultGetLinkWireSet,
)

func NewLinkCollectorModule(
	getLink *GetLinkHandler,
) *LinkCollectorModule {
	return &LinkCollectorModule{
		getLink: getLink,
	}
}

type LinkCollectorModule struct {
	getLink    *GetLinkHandler
	commandBus *cqrs.CommandBus
	eventBus   *cqrs.EventBus
}

func (l LinkCollectorModule) SetCommandBus(commandBus *cqrs.CommandBus) {
	l.commandBus = commandBus
}

func (l LinkCollectorModule) SetEventBus(eventBus *cqrs.EventBus) {
	l.eventBus = eventBus
}

func (l LinkCollectorModule) ListEventHandlers() []cqrs.EventHandler {
	return []cqrs.EventHandler{}
}

func (l LinkCollectorModule) ListCommandHandlers() []cqrs.CommandHandler {
	return []cqrs.CommandHandler{
		l.getLink,
	}
}
