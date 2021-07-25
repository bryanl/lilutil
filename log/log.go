package log

import (
	"context"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

type ctxValueType string

const (
	logKey ctxValueType = "_logger"
)

// From extracts a logger from a context. If one does not exist, a new logger is created.
func From(ctx context.Context) logr.Logger {
	if ctx == nil {
		return newLogger()
	}

	if logger, ok := ctx.Value(logKey).(logr.Logger); ok {
		return logger
	}

	return newLogger()
}

// WithLogger creates a new context with an embedded logger.
func WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey, newLogger())
}

func newLogger() logr.Logger {
	return logrusr.NewLogger(logrus.New())
}
