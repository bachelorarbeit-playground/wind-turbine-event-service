package config

import (
	"event-service/pkg/jetstream"

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
var Service = "EVENT_SERVICE"

func init() {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Debug().Err(err).Msg("Didn't found app.env file")
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
