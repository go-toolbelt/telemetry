package telemetry_test

import (
	"context"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/require"

	"github.com/go-toolbelt/telemetry"
)

func Test(t *testing.T) {
	ctx := context.Background()

	defer leaktest.Check(t)()

	instance := telemetry.New(telemetry.Options{
		Log: telemetry.LogOptions{
			Stdout: true,
		},
		Trace: telemetry.TraceOptions{
			Stdout: true,
		},
		Metric: telemetry.MetricOptions{
			Stdout: true,
		},
	})
	require.NoError(t, instance.Start())

	logger := instance.Logger("test")
	require.NotNil(t, logger)

	logger.Error("test message")

	tracer := instance.Tracer("test")
	require.NotNil(t, tracer)

	_, span := tracer.Start(ctx, "test span")
	span.End()

	meter := instance.Meter("test")
	require.NotNil(t, meter)

	counter, err := meter.NewInt64Counter("test counter")
	require.NoError(t, err)
	counter.Add(ctx, 1)

	require.NoError(t, instance.Shutdown(ctx))
}
