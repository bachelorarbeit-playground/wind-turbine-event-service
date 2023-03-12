package storage

import (
	"event-sourcing/pkg/jmespath"
	"event-sourcing/pkg/model"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	AverageProduction map[string]*AverageProductionEntry
	AnomalyDetection  map[string]*AnomalyDetectionEntry
}

var Store Storage

func init() {
	Store = Storage{
		AverageProduction: make(map[string]*AverageProductionEntry),
		AnomalyDetection:  make(map[string]*AnomalyDetectionEntry),
	}
}

func GetMaterializedViewEntries[T any](mv map[string]T) []T {
	entries := make([]T, 0)

	for _, v := range mv {
		entries = append(entries, v)
	}

	return entries
}

func (s *Storage) UpdateAverageProduction(event model.RawWindDataPayload) {

	key := event.ParkId + "_" + event.Region

	val, ok := s.AverageProduction[key]
	if !ok {
		s.AverageProduction[key] = &AverageProductionEntry{
			ParkId: event.ParkId,

			Region: event.Region,

			AverageValue:       event.Value,
			count_AverageValue: 1,
			sum_AverageValue:   event.Value,
			HighestValue:       event.Value,
			LowestValue:        event.Value,
		}
		return
	}

	val.count_AverageValue++
	val.sum_AverageValue += event.Value
	val.AverageValue = val.sum_AverageValue / float64(val.count_AverageValue)

	if event.Value > val.HighestValue {
		val.HighestValue = event.Value
	}

	if event.Value < val.LowestValue {
		val.LowestValue = event.Value
	}

}

func (s *Storage) UpdateAnomalyDetection(event model.RawWindDataPayload) {
	res, err := jmespath.TestCondition("availability <= `50`", event)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	if res {

		key := uuid.NewV4().String()

		s.AnomalyDetection[key] = &AnomalyDetectionEntry{
			ParkId: event.ParkId,

			TurbineId: event.TurbineId,

			Value: event.Value,

			Availability: event.Availability,
		}

	}

}
