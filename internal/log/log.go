package log

import (
	"errors"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var Logger zerolog.Logger

func InitializeLogger() error {
	if os.Getenv("ENV_TYPE") == "production" {
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
			With().
			Timestamp().
			Logger()

	} else {
		Logger = zerolog.New(os.Stdout).
			Output(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Logger()
	}

	return nil
}
