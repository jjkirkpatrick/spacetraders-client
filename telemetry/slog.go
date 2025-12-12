// Package telemetry provides OpenTelemetry integration utilities for the SpaceTraders client.
// This package exposes slog handlers that send logs to both console and OTLP endpoints.
//
// Example usage:
//
//	import (
//		"log/slog"
//		"os"
//
//		"github.com/jjkirkpatrick/spacetraders-client/client"
//		"github.com/jjkirkpatrick/spacetraders-client/telemetry"
//	)
//
//	func main() {
//		// Configure client with OTel
//		opts := client.DefaultClientOptions()
//		opts.TelemetryOptions = client.DefaultTelemetryOptions()
//		opts.TelemetryOptions.OTLPEndpoint = "localhost:4317"
//		c, _ := client.NewClient(opts)
//
//		// Set up combined logging (console + OTLP/Loki)
//		consoleHandler := slog.NewTextHandler(os.Stdout, nil)
//		handler := telemetry.NewCombinedSlogHandler("my-service", slog.LevelInfo, consoleHandler)
//		slog.SetDefault(slog.New(handler))
//	}
package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"
)

// OTelSlogHandler is an slog.Handler that sends logs to OpenTelemetry.
// It automatically attaches trace context (trace_id, span_id) to logs
// when called within a traced context.
type OTelSlogHandler struct {
	logger log.Logger
	attrs  []slog.Attr
	groups []string
	level  slog.Level
}

// NewOTelSlogHandler creates a new slog handler that exports logs via OpenTelemetry.
// The serviceName should match your telemetry service name for correlation.
// Logs will be sent to the OTLP endpoint configured in your client's TelemetryOptions.
func NewOTelSlogHandler(serviceName string, level slog.Level) *OTelSlogHandler {
	return &OTelSlogHandler{
		logger: global.GetLoggerProvider().Logger(serviceName),
		level:  level,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *OTelSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle processes the log record and sends it to OpenTelemetry.
func (h *OTelSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Convert slog level to OTel severity
	severity := slogLevelToOTelSeverity(record.Level)

	// Build the log record
	logRecord := log.Record{}
	logRecord.SetTimestamp(record.Time)
	logRecord.SetSeverity(severity)
	logRecord.SetSeverityText(record.Level.String())
	logRecord.SetBody(log.StringValue(record.Message))

	// Add trace context if available
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		logRecord.AddAttributes(
			log.String("trace_id", spanCtx.TraceID().String()),
			log.String("span_id", spanCtx.SpanID().String()),
		)
		if spanCtx.IsSampled() {
			logRecord.AddAttributes(log.Bool("trace_sampled", true))
		}
	}

	// Add pre-configured attributes
	for _, attr := range h.attrs {
		logRecord.AddAttributes(slogAttrToOTel(attr))
	}

	// Add record attributes
	record.Attrs(func(attr slog.Attr) bool {
		logRecord.AddAttributes(slogAttrToOTel(attr))
		return true
	})

	// Emit the log
	h.logger.Emit(ctx, logRecord)

	return nil
}

// WithAttrs returns a new handler with the given attributes added.
func (h *OTelSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &OTelSlogHandler{
		logger: h.logger,
		attrs:  make([]slog.Attr, len(h.attrs)+len(attrs)),
		groups: h.groups,
		level:  h.level,
	}
	copy(newHandler.attrs, h.attrs)
	copy(newHandler.attrs[len(h.attrs):], attrs)
	return newHandler
}

// WithGroup returns a new handler with the given group name.
func (h *OTelSlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newHandler := &OTelSlogHandler{
		logger: h.logger,
		attrs:  h.attrs,
		groups: append(h.groups, name),
		level:  h.level,
	}
	return newHandler
}

// slogLevelToOTelSeverity converts slog.Level to OTel log.Severity
func slogLevelToOTelSeverity(level slog.Level) log.Severity {
	switch {
	case level >= slog.LevelError:
		return log.SeverityError
	case level >= slog.LevelWarn:
		return log.SeverityWarn
	case level >= slog.LevelInfo:
		return log.SeverityInfo
	default:
		return log.SeverityDebug
	}
}

// slogAttrToOTel converts an slog.Attr to an OTel log.KeyValue
func slogAttrToOTel(attr slog.Attr) log.KeyValue {
	key := attr.Key
	value := attr.Value

	switch value.Kind() {
	case slog.KindString:
		return log.String(key, value.String())
	case slog.KindInt64:
		return log.Int64(key, value.Int64())
	case slog.KindUint64:
		return log.Int64(key, int64(value.Uint64()))
	case slog.KindFloat64:
		return log.Float64(key, value.Float64())
	case slog.KindBool:
		return log.Bool(key, value.Bool())
	case slog.KindTime:
		return log.String(key, value.Time().Format("2006-01-02T15:04:05.000Z07:00"))
	case slog.KindDuration:
		return log.String(key, value.Duration().String())
	case slog.KindGroup:
		// For groups, flatten with dot notation
		attrs := value.Group()
		if len(attrs) == 1 {
			return slogAttrToOTel(slog.Attr{Key: key + "." + attrs[0].Key, Value: attrs[0].Value})
		}
		// For multiple attrs in a group, just stringify
		return log.String(key, value.String())
	default:
		return log.String(key, value.String())
	}
}

// CombinedSlogHandler combines the OTel handler with a console handler
// so logs go to both stdout and OTel. This is the recommended handler
// for most applications.
type CombinedSlogHandler struct {
	otelHandler    *OTelSlogHandler
	consoleHandler slog.Handler
}

// NewCombinedSlogHandler creates a handler that logs to both console and OTel.
// This is the recommended way to set up logging with the SpaceTraders client.
//
// Parameters:
//   - serviceName: Should match your TelemetryOptions.ServiceName for log correlation
//   - level: Minimum log level to emit
//   - consoleHandler: Handler for console output (e.g., slog.NewTextHandler or slog.NewJSONHandler)
//
// Example:
//
//	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
//	handler := telemetry.NewCombinedSlogHandler("my-service", slog.LevelInfo, consoleHandler)
//	slog.SetDefault(slog.New(handler))
func NewCombinedSlogHandler(serviceName string, level slog.Level, consoleHandler slog.Handler) *CombinedSlogHandler {
	return &CombinedSlogHandler{
		otelHandler:    NewOTelSlogHandler(serviceName, level),
		consoleHandler: consoleHandler,
	}
}

// Enabled returns true if either handler is enabled.
func (h *CombinedSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.otelHandler.Enabled(ctx, level) || h.consoleHandler.Enabled(ctx, level)
}

// Handle sends the record to both handlers.
func (h *CombinedSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Send to console
	if h.consoleHandler.Enabled(ctx, record.Level) {
		h.consoleHandler.Handle(ctx, record)
	}
	// Send to OTel
	if h.otelHandler.Enabled(ctx, record.Level) {
		h.otelHandler.Handle(ctx, record)
	}
	return nil
}

// WithAttrs returns a new handler with attributes added to both handlers.
func (h *CombinedSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CombinedSlogHandler{
		otelHandler:    h.otelHandler.WithAttrs(attrs).(*OTelSlogHandler),
		consoleHandler: h.consoleHandler.WithAttrs(attrs),
	}
}

// WithGroup returns a new handler with a group added to both handlers.
func (h *CombinedSlogHandler) WithGroup(name string) slog.Handler {
	return &CombinedSlogHandler{
		otelHandler:    h.otelHandler.WithGroup(name).(*OTelSlogHandler),
		consoleHandler: h.consoleHandler.WithGroup(name),
	}
}
