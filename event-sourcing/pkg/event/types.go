package event

import (
	"event-sourcing/pkg/jetstream"
	"event-sourcing/pkg/model"
)

type EventStatus int

const (
	StatusPassed EventStatus = iota
	StatusDropped
	StatusFailed
)

type ForwardFunc[T any] func(model.RawWindDataPayload, jetstream.Subject) (EventStatus, error)

type State struct {
	Id       string `json:"id"`
	Instance string `json:"instance"`
	State    string `json:"state"`
	Error    string `json:"error"`
}
