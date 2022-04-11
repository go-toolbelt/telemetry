package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	log "go.uber.org/zap"
)

type Options struct {
	Log    LogOptions
	Trace  TraceOptions
	Metric MetricOptions
}

type Telemetry struct {
	options          Options
	tracerProvider   tracerProvider
	metricController metricController
	loggerConfig     log.Config
}

func New(options Options) *Telemetry {
	return &Telemetry{
		options: options,
	}
}

func (telemetry *Telemetry) Start() error {
	var err error

	telemetry.loggerConfig = newLogConfig(telemetry.options.Log)

	telemetry.tracerProvider, err = newTracerProvider(telemetry.options.Trace)
	if err != nil {
		return err
	}

	telemetry.metricController, err = newMetricController(telemetry.options.Metric)
	if err != nil {
		return err
	}

	return nil
}

func (telemetry *Telemetry) Shutdown(ctx context.Context) error {
	if err := telemetry.tracerProvider.Shutdown(ctx); err != nil {
		return err
	}

	if err := telemetry.metricController.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func (telemetry *Telemetry) Logger(name string, opts ...log.Option) *log.Logger {
	nameOpt := log.Fields(log.String("instrumentation.name", name))

	logger, err := telemetry.loggerConfig.Build(append(opts, nameOpt)...)
	if err != nil {
		return log.NewNop()
	}

	return logger
}

func (telemetry *Telemetry) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return telemetry.tracerProvider.Tracer(name, opts...)
}

func (telemetry *Telemetry) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return telemetry.metricController.Meter(name, opts...)
}
