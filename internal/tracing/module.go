package tracing

import (
	"context"

	"github.com/go-logr/zapr"
	"go.opentelemetry.io/otel"

	// "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"tracing",
		fx.Provide(
			func() (trace.SpanExporter, error) {
				// return stdouttrace.New(
				// 	stdouttrace.WithPrettyPrint(),
				//	stdouttrace.WithoutTimestamps(),
				// )

				return jaeger.New(
					jaeger.WithAgentEndpoint(jaeger.WithAgentHost("jaeger")),
				)
			},
			func(
				exporter trace.SpanExporter,
				resource *resource.Resource,
			) *trace.TracerProvider {
				return trace.NewTracerProvider(
					trace.WithBatcher(exporter),
					trace.WithResource(resource),
				)
			},
		),
		fx.Invoke(
			func(
				lc fx.Lifecycle,
				logger *zap.Logger,
				traceProvider *trace.TracerProvider,
			) {
				otel.SetLogger(zapr.NewLogger(logger))

				lc.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							return traceProvider.Shutdown(ctx)
						},
					},
				)

				otel.SetTracerProvider(traceProvider)
			},
		),
	)
}
