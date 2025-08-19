package log

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() (*zerolog.Logger, error) {
	if os.Getenv("ENV_TYPE") == "development" {
		logger := zerolog.New(os.Stdout).
			Output(zerolog.ConsoleWriter{Out: os.Stdout}).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Logger()

		return &logger, nil
	}

	logDirectory := os.Getenv("LOG_DIRECTORY")

	if len(logDirectory) == 0 {
		return nil, errors.New("LOG_DIRECTORY is not set")
	}

	rotationLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDirectory, "app.log"),
		MaxSize:    500,
		MaxBackups: 10,
		MaxAge:     30,
	}
	logger := zerolog.New(rotationLogger).
		Output(zerolog.ConsoleWriter{Out: rotationLogger}).
		Level(zerolog.ErrorLevel).
		With().
		Timestamp().
		Logger()

	return &logger, nil
}
