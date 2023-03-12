package event

import (
	"encoding/json"
	"fmt"

	"event-sourcing/pkg/config"
	"event-sourcing/pkg/jetstream"
	"event-sourcing/pkg/lua"
	"event-sourcing/pkg/model"
	"event-sourcing/pkg/storage"

	"github.com/adimunteanu/gluamapper"
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

	var forwardMap = map[jetstream.Subject]ForwardFunc[model.Payload]{

		"ingestionPipeline": forwardEvent[*model.ProcessedDataPayload],

		"anomalyDetection": forwardEvent[*model.AnomalyDataPayload],
	}

	for subject, forwardFn := range forwardMap {
		status, err := forwardFn(inputEvent, subject)
		if err != nil {
			log.Warn().Str("subject", string(subject)).Str("tracing", tracingID).Msg(err.Error())
		}

		if status != StatusDropped {
			sendStatusEvent(err, tracingID)
		}
	}
}

func sendEvent(msgData []byte, subject jetstream.Subject) error {
	// create connection to NATS server
	nc, err := nats.Connect(config.ServiceConfig.NatsServer)
	if err != nil {
		return fmt.Errorf("connection to NATS server failed: %s", err)
	}
	defer nc.Close()

	// send events
	msg := nats.Msg{Subject: string(subject), Data: msgData}
	nc.PublishMsg(&msg)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		return fmt.Errorf("sending processed events failed: %s", err)
	}

	return nil
}

func forwardEvent[T model.Payload](event model.RawWindDataPayload, subject jetstream.Subject) (EventStatus, error) {
	output, err := lua.RunScript(config.AppConfig.Script, string(subject), event)

	if err != nil {
		return StatusFailed, fmt.Errorf("issue running script: %s", err)
	}

	if output == nil {
		return StatusDropped, nil
	}

	options := gluamapper.Option{WeaklyTypedInput: false, ErrorUnused: true, ErrorUnset: true}
	mapper := gluamapper.NewMapper(options)

	var processedEvent T
	err = mapper.Map(output, &processedEvent)
	if err != nil {
		return StatusFailed, fmt.Errorf("failed to cast a lua table: %s", err)
	}

	msgData, err := processedEvent.Marshal()
	if err != nil {
		return StatusFailed, fmt.Errorf("could not marshal output of script: %s", err)
	}

	err = sendEvent(msgData, subject)
	if err != nil {
		return StatusFailed, err
	}

	return StatusPassed, nil
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
