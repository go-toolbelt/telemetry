package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type MetricOptions struct {
	Stdout bool
}

type metricController interface {
	Meter(instrumentationName string, opts ...metric.MeterOption) metric.Meter
	Stop(ctx context.Context) error
}

func newMetricController(options MetricOptions) (metricController, error) {
	var metricExporter export.Exporter

	if options.Stdout {
		var err error

		metricExporter, err = stdoutmetric.New()
		if err != nil {
			return nil, err
		}
	}

	metricController := controller.New(
		newCheckpointerFactory(metricExporter),
		controller.WithExporter(metricExporter),
	)

	if err := metricController.Start(context.Background()); err != nil {
		return nil, err
	}

	return metricController, nil
}

type checkpointerFactory struct {
	metricExporter export.Exporter
}

func newCheckpointerFactory(metricExporter export.Exporter) export.CheckpointerFactory {
	return &checkpointerFactory{
		metricExporter: metricExporter,
	}
}

func (factory *checkpointerFactory) NewCheckpointer() export.Checkpointer {
	return processor.New(
		selector.NewWithHistogramDistribution(),
		factory.metricExporter,
	)
}
