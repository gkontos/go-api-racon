package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	Logger zerolog.Logger
)

func initModule(debug, console_log bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if console_log {
		Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		log.Logger = log.Output(os.Stdout).With().Timestamp().Logger()
		Logger = log.Logger
	}

}

func InitLogger(debug, console_log bool) {
	initModule(debug, console_log)
}
