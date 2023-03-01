package modules

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/aiocean/shopify-crawler/internal/modules/linkcollector"
	"github.com/google/wire"
)

type Module interface {
	ListEventHandlers() []cqrs.EventHandler
	ListCommandHandlers() []cqrs.CommandHandler
}

var DefaultModuleWireSet = wire.NewSet(
	NewModuleList,
	linkcollector.DefaultLinkCollectorWireSet,
)

func NewModuleList(
	linkCollector *linkcollector.LinkCollectorModule,
) []*cqrs.CommandHandler {
	return []Module{
		linkCollector,
	}
}
