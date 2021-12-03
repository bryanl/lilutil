package dbutil

import (
	"context"

	"github.com/bryanl/lilutil/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
