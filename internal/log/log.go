package log

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

func InitializeLogger() error {
	if os.Getenv("ENV_TYPE") == "development" {
		Logger = zerolog.New(os.Stdout).
			Output(zerolog.ConsoleWriter{Out: os.Stdout}).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Logger()

		return nil
	}

	logDirectory := os.Getenv("LOG_DIRECTORY")

	if len(logDirectory) == 0 {
		return errors.New("LOG_DIRECTORY is not set")
	}

	rotationLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logDirectory, "app.log"),
		MaxSize:    500,
		MaxBackups: 10,
		MaxAge:     30,
	}
	Logger = zerolog.New(rotationLogger).
		Output(zerolog.ConsoleWriter{Out: rotationLogger}).
		Level(zerolog.ErrorLevel).
		With().
		Timestamp().
		Logger()

	return nil
}
