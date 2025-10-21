package logger

import (
	appconfig "logger-example/config"

	"go.uber.org/fx"
)

// FX Module for logger
var Module = fx.Module("logger",
	fx.Provide(NewLoggerFromConfig),
)

// NewLoggerFromConfig creates a logger from config
func NewLoggerFromConfig(cfg *appconfig.Config) {
	Init(Config{
		Level:           cfg.Logger.Level,
		Encoding:        cfg.Logger.Encoding,
		Service:         cfg.Logger.Service,
		Environment:     cfg.Logger.Environment,
		AsyncBufferSize: cfg.Logger.AsyncBufferSize,
		BatchSize:       cfg.Logger.BatchSize,
	})
}
