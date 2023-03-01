//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/aiocean/shopify-crawler/internal/cqrssvc"
	"github.com/aiocean/shopify-crawler/internal/logsvc"
	modules "github.com/aiocean/shopify-crawler/internal/modules"
	"github.com/aiocean/shopify-crawler/internal/server"
	"github.com/google/wire"
)

func InitializeHandler(ctx context.Context) (*server.Server, func(), error) {
	wire.Build(
		server.DefaultServerWireSet,
		modules.DefaultModuleWireSet,
		cqrssvc.DefaultCqrsWireSet,
		logsvc.DefaultWireSet,
	)

	return nil, nil, nil
}
