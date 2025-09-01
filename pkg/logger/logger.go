package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
		NoColor:    false,
	}

	Logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

func Info(msg string) {
	Logger.Info().Msg(msg)
}

func Fatal(err error) {
	Logger.Fatal().Err(err).Msg("Fatal error occurred")
}

func Error(err error) {
	Logger.Error().Err(err).Msg("Error occurred")
}

func Warn(msg string) {
	Logger.Warn().Msg(msg)
}
