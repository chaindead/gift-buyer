package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}

	log.Logger = log.Output(writer).With().Caller().Stack().Logger().Level(zerolog.DebugLevel)
}
