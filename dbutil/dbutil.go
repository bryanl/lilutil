package dbutil

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bryanl/lilutil/log"
)

// InitSqlite initializes a sqlite db connection.
func InitSqlite(ctx context.Context, dsn string, models ...interface{}) (*gorm.DB, error) {
	logger := log.From(ctx)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: log.DBLogger(logger),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(models...); err != nil {
		return nil, err
	}

	return db, nil
}

// InitPostgres initialize a postgres db connection.
func InitPostgres(ctx context.Context, dsn string, models ...interface{}) (*gorm.DB, error) {
	logger := log.From(ctx)

	logger.Info("Initializing postgres database")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log.DBLogger(logger),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(models...); err != nil {
		return nil, err
	}

	logger.Info("Database initialized")

	return db, nil
}
