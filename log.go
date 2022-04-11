package telemetry

import (
	"go.uber.org/zap"
)

type LogOptions struct {
	Zap zap.Config
}

func newLogConfig(options LogOptions) zap.Config {
	return options.Zap
}
