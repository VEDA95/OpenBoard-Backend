package db

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormZerologger struct {
	Logger               *zerolog.Logger
	SlowThreshold        time.Duration
	SkipCallerLookup     bool
	IgnoreRecordNotFound bool
}

func NewGormZerologger(logger *zerolog.Logger) *GormZerologger {
	return &GormZerologger{
		Logger:               logger,
		SlowThreshold:        200 * time.Millisecond,
		SkipCallerLookup:     false,
		IgnoreRecordNotFound: true,
	}
}

func (l *GormZerologger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	switch level {
	case gormlogger.Silent:
		zeroLevel := zerolog.Disabled
		l.Logger.Level(zeroLevel)
	case gormlogger.Error:
		l.Logger.Level(zerolog.ErrorLevel)
	case gormlogger.Warn:
		l.Logger.Level(zerolog.WarnLevel)
	case gormlogger.Info:
		l.Logger.Level(zerolog.InfoLevel)
	}

	return l
}

func (logger *GormZerologger) Info(ctx context.Context, msg string, data ...any) {
	logger.Logger.Info().Msgf(msg, data...)
}

func (logger *GormZerologger) Warn(ctx context.Context, msg string, data ...any) {
	logger.Logger.Warn().Msgf(msg, data...)
}

func (logger *GormZerologger) Error(ctx context.Context, msg string, data ...any) {
	logger.Logger.Error().Msgf(msg, data...)
}

func (logger *GormZerologger) Trace(ctx context.Context, begin time.Time, callback func() (sql string, rowsAffected int64), err error) {
	if logger.Logger == nil {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := callback()

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !logger.IgnoreRecordNotFound):
		logger.Logger.Error().
			Err(err).
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msg("Database query error")

	case elapsed > logger.SlowThreshold && logger.SlowThreshold != 0:
		logger.Logger.Warn().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msgf("Slow SQL query [%v]", elapsed)

	case logger.Logger.GetLevel() <= zerolog.DebugLevel:
		logger.Logger.Debug().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msg("Database query executed")
	}
}
