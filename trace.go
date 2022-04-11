package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type tracerProvider interface {
	trace.TracerProvider

	Shutdown(ctx context.Context) error
}

type TraceOptions struct {
	Stdout bool
}

func newTracerProvider(options TraceOptions) (tracerProvider, error) {
	var providerOptions []sdk.TracerProviderOption

	if options.Stdout {
		exporter, err := stdouttrace.New()
		if err != nil {
			return nil, err
		}

		providerOptions = append(providerOptions, sdk.WithSpanProcessor(
			sdk.NewBatchSpanProcessor(
				exporter,
			),
		))
	}

	if len(providerOptions) == 0 {
		return newNoopTracerProvider(), nil
	}

	return sdk.NewTracerProvider(providerOptions...), nil
}

type noopTracerProvider struct {
	trace.TracerProvider
}

func newNoopTracerProvider() tracerProvider {
	return noopTracerProvider{
		TracerProvider: trace.NewNoopTracerProvider(),
	}
}

func (noopTracerProvider) Shutdown(ctx context.Context) error {
	return nil
}
