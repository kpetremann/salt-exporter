package logging

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure load log default configuration (like format, output target etc...).
func Configure() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// SetLogLevel configures the loglevel.
//
// logLevel: The log level to use, in zerolog format.
func SetLevel(logLevel string) {
	level, err := zerolog.ParseLevel(logLevel)
	fmt.Println(logLevel)
	if err != nil {
		fmt.Println("Failed to parse log level")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(level)
}
