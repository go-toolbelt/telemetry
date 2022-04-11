package telemetry

import (
	log "go.uber.org/zap"
)

type LogOptions struct {
	Stdout bool
}

func newLogConfig(_ LogOptions) log.Config {
	return log.NewProductionConfig()
}
