package main

import (
	"io/ioutil"
	"os"

	"event-sourcing/pkg/config"
	"event-sourcing/pkg/event"
	"event-sourcing/pkg/file"
	"event-sourcing/pkg/handler"
	"event-sourcing/pkg/jetstream"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set up zerolog time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Set pretty logging on
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if !file.CheckIfFileExists("processing/processing.lua") {
		log.Fatal().Msg("script file not found")
	}

	script, err := ioutil.ReadFile("processing/processing.lua")

	if err != nil {
		log.Fatal().Err(err).Msg("Could not read script from file")
	}

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
