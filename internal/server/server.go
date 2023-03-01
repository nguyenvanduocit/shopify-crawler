package server

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/aiocean/shopify-crawler/internal/cqrssvc"
	"github.com/aiocean/shopify-crawler/internal/modules"
	"github.com/google/wire"
)

var DefaultServerWireSet = wire.NewSet(NewServer)

type Server struct {
	Done          chan error
	messageRouter *message.Router
}

func NewServer(
	modules []modules.Module,
	cqrsSvc *cqrssvc.CqrsSvc,
	messageRouter *message.Router,

) (*Server, func(), error) {

	for _, module := range modules {
		commands := module.ListCommandHandlers()
		for _, command := range commands {
			cqrsSvc.AddCommandHandler(command)
		}

		events := module.ListEventHandlers()
		for _, event := range events {
			cqrsSvc.AddEventHandler(event)
		}
	}

	cleanup := func() {
		messageRouter.Close()
	}

	return &Server{
		messageRouter: messageRouter,
		Done:          make(chan error),
	}, cleanup, nil
}

func (s *Server) Start() {
	err := s.messageRouter.Run(context.Background())
	if err != nil {
		s.Done <- err
		return
	}

	s.Done <- nil
}
