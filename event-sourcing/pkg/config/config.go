package config

import (
	"event-sourcing/pkg/jetstream"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type serviceConfig struct {
	Name       string
	NatsServer string
}

type appConfig struct {
	Script         []byte
	InputStream    jetstream.Stream
	OutputSubjects []jetstream.Subject
}

var AppConfig appConfig
var ServiceConfig serviceConfig
var Service = "event-sourcing"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Debug().Err(err).Msg("Didn't found .env file")
	}

	ServiceConfig = serviceConfig{}
	err = envconfig.Process(Service, &ServiceConfig)
	if err != nil {
		log.Fatal().Msg("Could not parse config:")
	}

	AppConfig = appConfig{
		InputStream: jetstream.Stream{
			Name:    "EventLog",
			Subject: "rawWindData",
		},
		OutputSubjects: []jetstream.Subject{

			"ingestionPipeline",

			"anomalyDetection",
		},
	}
}
