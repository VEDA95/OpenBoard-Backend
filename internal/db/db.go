package db

import (
	"VEDA95/open_board/api/internal/log"
	"errors"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	envType := os.Getenv("ENV_TYPE")

	if len(dsn) == 0 {
		return nil, errors.New("DATABASE_URL not set")
	}

	gormLogger := NewGormZerologger(log.Global)

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

	if envType != "production" {
		instacne = instacne.Debug()
	}

	sqlDB, err := instacne.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return instacne, nil
}
