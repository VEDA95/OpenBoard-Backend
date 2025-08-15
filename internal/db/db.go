package db

import (
	applogger "VEDA95/open_board/api/internal/log"
	"errors"
	"os"
	"reflect"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Instance *gorm.DB

func InitializeDB(models []any) error {
	dsn := os.Getenv("DATABASE_URL")

	if len(dsn) == 0 {
		return errors.New("DATABASE_URL not set")
	}

	if reflect.ValueOf(applogger.Logger).IsZero() {
		return errors.New("logger not set")
	}

	gormLogger := NewGormZerologger(&applogger.Logger)

	var logLevel logger.LogLevel
	if os.Getenv("ENV_TYPE") == "production" {
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
		return err
	}

	Instance = instacne

	sqlDB, err := Instance.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if os.Getenv("ENV") == "development" {
		if err := AutoMigrate(models); err != nil {
			return err
		}
	}

	return nil
}

func AutoMigrate(models []any) error {
	return Instance.AutoMigrate(models...)
}
