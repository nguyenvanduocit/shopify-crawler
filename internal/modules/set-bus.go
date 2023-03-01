package modules

import "github.com/ThreeDotsLabs/watermill/components/cqrs"

type SetBusable interface {
	SetCommandBus(commandBus *cqrs.CommandBus)
	SetEventBus(eventBus *cqrs.EventBus)
}
