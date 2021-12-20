package log

import (
	"context"
	"io"
	"os"

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

// LoggerOption is an option for configuring the logger.
type LoggerOption func(config *LoggerConfig)

// LoggerOutput sets the output location for the logger.
func LoggerOutput(w io.Writer) LoggerOption {
	if w == nil {
		panic("logger output cannot be nil")
	}

	return func(config *LoggerConfig) {
		config.out = w
	}
}

// WithExistingLogger creates a new context with an existing logger.
func WithExistingLogger(ctx context.Context, logger logr.Logger) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

// WithLogger creates a new context with an embedded logger.
func WithLogger(ctx context.Context, options ...LoggerOption) context.Context {
	return context.WithValue(ctx, logKey, newLogger(options...))
}

// LoggerConfig is logger configuration.
type LoggerConfig struct {
	out io.Writer
}

func newLoggerConfig() *LoggerConfig {
	config := &LoggerConfig{
		out: os.Stderr,
	}

	return config
}

func (config *LoggerConfig) update(logger *logrus.Logger) {
	logger.Out = config.out
}

func newLogger(options ...LoggerOption) logr.Logger {
	config := newLoggerConfig()
	for _, option := range options {
		option(config)
	}

	l := logrus.New()
	config.update(l)

	return logrusr.NewLogger(l)
}
