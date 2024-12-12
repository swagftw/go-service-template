package applog

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"go-service-template/internal/apperr"
)

var Logger ILogger

type ILogger interface {
	Error(ctx context.Context, err error, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Debug(ctx context.Context, msg string, fields map[string]interface{})
}

type logger struct {
	otelLogger *otelzap.Logger
}

var _ ILogger = logger{}

// InitLogger initializes the logger
func InitLogger(debug bool) error {
	// adds caller, skip the caller stack by 1 so it shows the caller of the error log, and adds stacktrace if error level
	zapLogger, err := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))

	if debug {
		zapLogger, err = zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	}

	if err != nil {
		slog.Error("failed to initialize logger", "err", err)

		return apperr.New(http.StatusInternalServerError, err, "failed to initialize logger", apperr.ErrInternalError)
	}

	otelLogger := otelzap.New(zapLogger)

	Logger = &logger{
		otelLogger: otelLogger,
	}

	return nil
}

func (l logger) Error(ctx context.Context, err error, msg string, fields map[string]interface{}) {
	f := make([]zap.Field, 0, len(fields))

	for k, v := range fields {
		f = append(f, zap.Any(k, v))
	}

	l.otelLogger.ErrorContext(ctx, msg, f...)
}

func (l logger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	f := make([]zap.Field, 0, len(fields))

	for k, v := range fields {
		f = append(f, zap.Any(k, v))
	}

	l.otelLogger.WarnContext(ctx, msg, f...)
}

func (l logger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	f := make([]zap.Field, 0, len(fields))

	for k, v := range fields {
		f = append(f, zap.Any(k, v))
	}

	l.otelLogger.InfoContext(ctx, msg, f...)
}

func (l logger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	f := make([]zap.Field, 0, len(fields))

	for k, v := range fields {
		f = append(f, zap.Any(k, v))
	}
	l.otelLogger.DebugContext(ctx, msg, f...)
}
