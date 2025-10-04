package setup

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

// InitTracing installs a no-op tracer provider placeholder so instrumentation can be added later.
func InitTracing(logger *log.Logger) func(context.Context) error {
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	logger.Println("telemetry tracing initialized (POC)")
	return tp.Shutdown
}
