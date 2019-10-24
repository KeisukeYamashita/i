package logging

import (
	"context"

	"github.com/go-logr/logr"
)

type ctxKey struct{}

// WithContext add context value to the logger
func WithContext(ctx context.Context, log logr.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

// FromContext retrieves logger from context value
func FromContext(ctx context.Context) logr.Logger {
	return ctx.Value(ctxKey{}).(logr.Logger)
}
