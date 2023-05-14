package logging

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogging(logLevel string) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		fmt.Println("Failed to parse log level")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
