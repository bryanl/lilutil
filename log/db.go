package log

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"gorm.io/gorm/logger"
)

// DB is a database logger.
type DB struct {
	logger  logr.Logger
	level   logger.LogLevel
	message string
}

// DBOption is a function for configuring DB.
type DBOption func(*DB)

// WithDBLevel configures the DB log level.
func WithDBLevel(level logger.LogLevel) DBOption {
	return func(d *DB) {
		d.level = level
	}
}

// WithDBLogMessage sets the log message for DB logs.
func WithDBLogMessage(msg string) DBOption {
	return func(d *DB) {
		d.message = msg
	}
}

// DBLogger creates an instance of DB>
func DBLogger(l logr.Logger, options ...DBOption) *DB {
	db := &DB{
		message: "db command",
	}

	for _, option := range options {
		option(db)
	}

	db.logger = l.WithName("db")

	return db
}

// LogMode sets the level.
func (d *DB) LogMode(level logger.LogLevel) logger.Interface {
	return DBLogger(d.logger, WithDBLevel(level))
}

// Info logs at the info level.
func (d *DB) Info(_ context.Context, s string, i ...interface{}) {
	d.logger.Info(fmt.Sprintf(s, i...))
}

// Warn logs at the warning level.
func (d *DB) Warn(_ context.Context, s string, i ...interface{}) {
	d.logger.Info(fmt.Sprintf(s, i...))
}

// Error logs at the error level.
func (d *DB) Error(_ context.Context, s string, i ...interface{}) {
	d.logger.Error(errors.New("db"), fmt.Sprintf(s, i...))
}

// Trace logs at the trace level.
func (d *DB) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []interface{}{
		"sql", sql,
		"elapsed", float64(elapsed.Nanoseconds()) / 1e6,
	}
	if rows > 0 {
		fields = append(fields, "rows", rows)
	}

	switch {
	case err != nil && (!errors.Is(err, logger.ErrRecordNotFound)):
		d.logger.V(9).Info(d.message, append([]interface{}{"db-level", "trace"}, fields...)...)
	default:
		d.logger.V(9).Info(d.message, append([]interface{}{"db-level", "info"}, fields...)...)
	}
}
