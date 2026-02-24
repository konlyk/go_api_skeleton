package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	fallbackLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if err := execute(); err != nil {
		fallbackLogger.Fatal().Err(err).Msg("execute command")
	}
}
