package tracing

import (
	appconfig "logger-example/config"

	"go.uber.org/fx"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Tracer struct {
	serviceName string
	environment string
}

func NewTracer(cfg *appconfig.Config) *Tracer {
	return &Tracer{
		serviceName: cfg.DataDog.ServiceName,
		environment: cfg.DataDog.Environment,
	}
}

func (t *Tracer) Start() {
	tracer.Start(
		tracer.WithServiceName(t.serviceName),
		tracer.WithEnv(t.environment),
	)
}

func (t *Tracer) Stop() {
	tracer.Stop()
}

// FX Module for tracing
var Module = fx.Module("tracing",
	fx.Provide(NewTracer),
)
