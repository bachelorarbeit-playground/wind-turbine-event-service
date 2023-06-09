package main

import (
	"gql-service/graph"
	"gql-service/pkg/config"
	"gql-service/pkg/event"
	"gql-service/pkg/jetstream"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Set up zerolog time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Set pretty logging on
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal().Msgf("%w", http.ListenAndServe(":"+port, nil))
}
