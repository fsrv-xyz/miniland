package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"ref.ci/fsrvcorp/miniland/userland/pkg/web"
)

func init() {
	zlog.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger().With().Caller().Logger()
}

func main() {
	address := os.Getenv("WEB_ADDRESS")
	if address == "" {
		address = ":8080"
	}
	zlog.Info().Msgf("Starting web server on %#v", address)
	web.Start(address)
}
