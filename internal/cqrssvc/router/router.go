package router

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/garsue/watermillzap"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var DefaultRouterWireSet = wire.NewSet(
	NewRouter,
)

func NewRouter(
	zapLogger *zap.Logger,
) (*message.Router, func(), error) {
	logger := watermillzap.NewLogger(zapLogger)
	router, err := message.NewRouter(message.RouterConfig{}, logger)

	if err != nil {
		return nil, nil, err
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		Retry{
			MaxRetries:      2,
			InitialInterval: time.Second * 1,
			Logger:          logger,
			OnFailed: func(msg *message.Message, err error) ([]*message.Message, error) {
				zapLogger.Error("error processing message", zap.Error(err))
				return nil, nil
			},
		}.Middleware,
	)

	cleanup := func() {
		if err := router.Close(); err != nil {
			zapLogger.Error("error closing router", zap.Error(err))
		}
	}

	return router, cleanup, nil
}
