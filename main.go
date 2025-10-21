package main

import (
	"context"
	"logger-example/config"
	"logger-example/controller"
	"logger-example/logger"
	"logger-example/server"
	"logger-example/tracing"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Provide configuration
		fx.Provide(config.GetInstance),

		// Include modules
		config.Module,
		logger.Module,
		tracing.Module,
		controller.Module,
		server.Module,

		// Lifecycle hooks
		fx.Invoke(startTracing),
		fx.Invoke(startServer),
	).Run()
}

// startTracing starts the DataDog tracer
func startTracing(tracer *tracing.Tracer, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			tracer.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			tracer.Stop()
			return nil
		},
	})
}

// startServer starts the HTTP server
func startServer(srv *server.Server, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Server will be stopped by the context cancellation
			return nil
		},
	})
}
