package logging

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
