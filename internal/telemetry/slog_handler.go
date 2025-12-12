package telemetry

import (
	"log/slog"

	publictelemetry "github.com/jjkirkpatrick/spacetraders-client/telemetry"
)

// OTelSlogHandler re-exports the public OTelSlogHandler for internal use.
type OTelSlogHandler = publictelemetry.OTelSlogHandler

// CombinedSlogHandler re-exports the public CombinedSlogHandler for internal use.
type CombinedSlogHandler = publictelemetry.CombinedSlogHandler

// NewOTelSlogHandler creates a new slog handler that exports logs via OpenTelemetry.
// This is a re-export of the public telemetry package function.
func NewOTelSlogHandler(serviceName string, level slog.Level) *OTelSlogHandler {
	return publictelemetry.NewOTelSlogHandler(serviceName, level)
}

// NewCombinedSlogHandler creates a handler that logs to both console and OTel.
// This is a re-export of the public telemetry package function.
func NewCombinedSlogHandler(serviceName string, level slog.Level, consoleHandler slog.Handler) *CombinedSlogHandler {
	return publictelemetry.NewCombinedSlogHandler(serviceName, level, consoleHandler)
}
