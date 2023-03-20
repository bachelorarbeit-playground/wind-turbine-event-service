package jetstream

import (
	"github.com/nats-io/nats.go"

	"github.com/rs/zerolog/log"
)

type JetStreamClient struct {
	js nats.JetStreamContext
}

func Connect(server string) (*nats.Conn, error) {
	nc, err := nats.Connect(server)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("connected successfully to nats jetstream cluster")
	return nc, nil
}

func NewClient(nc *nats.Conn) (*JetStreamClient, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	log.Info().Msg("jetStreamContext created successfully")
	return &JetStreamClient{js}, nil
}

func (c *JetStreamClient) CreateInputStream(stream Stream) error {
	log.Info().Msgf("creating stream %q and subject %q\n", stream.Name, stream.Subject)
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:     stream.Name,
		Subjects: []string{string(stream.Subject)},
	})
	if err != nil {
		return err
	}

	log.Info().Msg("input stream created successfully")
	return nil
}

func (c *JetStreamClient) CheckStream(stream Stream) error {
	_, err := c.js.StreamInfo(stream.Name)
	if err != nil {
		return err
	}

	log.Info().Msg("input stream exists")
	return nil
}

func (c *JetStreamClient) Subscribe(handler nats.MsgHandler, stream Stream) error {
	_, err := c.js.Subscribe(string(stream.Subject), handler)
	if err != nil {
		return err
	}

	log.Info().Msg("subscribed to input subject successfully")
	return nil
}
