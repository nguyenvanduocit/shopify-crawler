package linkcollector

import (
	"context"
	"fmt"
	"github.com/google/wire"
)

var DefaultGetLinkWireSet = wire.NewSet(NewGetLinkHandler)

func NewGetLinkHandler() (*GetLinkHandler, func(), error) {
	return &GetLinkHandler{}, func() {}, nil
}

type GetLinkHandler struct {
}

type GetLinkCmd struct {
}

func (h *GetLinkHandler) HandlerName() string {
	return "GetLink"
}

func (h *GetLinkHandler) NewCommand() interface{} {
	return &GetLinkCmd{}
}

func (h *GetLinkHandler) Handle(ctx context.Context, cmdItf interface{}) error {
	cmd := cmdItf.(*GetLinkCmd)
	fmt.Println("GetLinkCmd", cmd)
	return nil
}
