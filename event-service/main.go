package main

import (
	_ "embed"
	"os"

	"event-service/pkg/config"
	"event-service/pkg/event"
	"event-service/pkg/handler"
	"event-service/pkg/jetstream"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed processing/processing.lua
var script []byte

func main() {
	// Set up zerolog time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Set pretty logging on
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config.AppConfig.Script = script

	nc, err := jetstream.Connect(config.ServiceConfig.NatsServer)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to nats jetstream cluster")
	}
	defer nc.Close()

	jsc, err := jetstream.NewClient(nc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create JetStreamContext")
	}

	err = jsc.CheckStream(config.AppConfig.InputStream)
	if err != nil {
		log.Warn().Msg("input stream doesn't exist")

		err = jsc.CreateInputStream(config.AppConfig.InputStream)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create input stream")
		}
	}

	err = jsc.Subscribe(event.HandleEvent, config.AppConfig.InputStream)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to subscribe to input subject")
	}

	e := echo.New()
	handler.NewHandler(e)
	e.Logger.Fatal(e.Start(":1323"))
}
