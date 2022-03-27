package logging

import (
	"context"

	"github.com/itiky/alienInvasion/pkg"
	"github.com/rs/zerolog"
)

const (
	// contextKeyLogger defines the key for the logger to be stored within request context.
	contextKeyLogger = pkg.ContextKey("Logger")
)

// GetCtxLogger returns a logger stored within the context or creates a new instance.
func GetCtxLogger(ctx context.Context) (context.Context, zerolog.Logger) {
	if ctx == nil {
		ctx = context.Background()
	}

	if ctxValue := ctx.Value(contextKeyLogger); ctxValue != nil {
		if ctxLogger, ok := ctxValue.(zerolog.Logger); ok {
			return ctx, ctxLogger
		}
	}
	logger := NewLogger()

	return SetCtxLogger(ctx, logger), logger
}

// SetCtxLogger adds the logger to the context overwriting an existing one.
func SetCtxLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}
