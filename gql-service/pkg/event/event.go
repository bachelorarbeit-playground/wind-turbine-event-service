package event

import (
	"encoding/json"

	"gql-service/pkg/config"
	"gql-service/pkg/model"
	"gql-service/pkg/storage"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func HandleEvent(m *nats.Msg) {
	tracingID := m.Header.Get("Tracing")
	if tracingID == "" {
		tracingID = uuid.NewV4().String()
	}
	log.Info().Str("subject", m.Subject).Str("tracing", tracingID).Msg("Received message")

	m.Ack()

	var inputEvent model.RawWindDataPayload
	err := inputEvent.Unmarshal(m.Data)

	// validation of input message failed
	if err != nil {
		log.Warn().Str("subject", m.Subject).Str("tracing", tracingID).Msg(err.Error())
		sendStatusEvent(err, tracingID)
		return
	}

	storage.Store.UpdateAverageProduction(inputEvent)

	storage.Store.UpdateAnomalyDetection(inputEvent)

}

func sendStatusEvent(err error, tracingID string) {
	data := State{
		Id:       tracingID,
		Instance: config.ServiceConfig.Name,
		State:    "passed",
		Error:    "",
	}
	if err != nil {
		// failed event
		data.State = "failed"
		data.Error = err.Error()
	}

	// event data in JSON
	msgData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Str("tracing", tracingID).Msg("could not marshal state")
		return
	}

	// create connection to NATS server
	nc, err := nats.Connect(config.ServiceConfig.NatsServer)
	if err != nil {
		log.Error().Err(err).Str("url", config.ServiceConfig.NatsServer).Str("tracing", tracingID).Msg("connection to NATS server failed")
	}
	defer nc.Close()

	// send events
	msg := nats.Msg{
		Subject: "state." + config.ServiceConfig.Name,
		Header:  nats.Header{"Tracing": []string{tracingID}},
		Data:    msgData,
	}
	nc.PublishMsg(&msg)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Error().Err(err).Str("url", config.ServiceConfig.NatsServer).Str("tracing", tracingID).Msg("publishing state to NATS server failed")
		return
	}
}
