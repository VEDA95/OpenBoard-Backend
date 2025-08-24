package db

import (
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(appLogger *zerolog.Logger, models []any) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	envType := os.Getenv("ENV_TYPE")

	if len(dsn) == 0 {
		return nil, errors.New("DATABASE_URL not set")
	}

	gormLogger := NewGormZerologger(appLogger)

	var logLevel logger.LogLevel
	if envType == "production" {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	instacne, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := instacne.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if envType == "development" {
		if err := instacne.AutoMigrate(models...); err != nil {
			return nil, err
		}
	}

	return instacne, nil
}

func AutoMigrate(dbInstance *gorm.DB, models []any) error {
	return dbInstance.AutoMigrate(models...)
}
