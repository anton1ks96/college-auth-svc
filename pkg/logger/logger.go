package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var newLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
		NoColor:    false,
	}

	newLogger = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func Debug(msg string) {
	newLogger.Debug().Msg(msg)
}

func Info(msg string) {
	newLogger.Info().Msg(msg)
}

func Fatal(err error) {
	newLogger.Fatal().Err(err)
}

func Error(err error) {
	newLogger.Error().Err(err)
}
