package runtime

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func (r *Runtime) logger() *Runtime {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}).With().Timestamp().Caller().Logger()

	r.Logger = logger

	return r
}
