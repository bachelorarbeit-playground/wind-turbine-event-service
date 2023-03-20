package model

import (
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"

	"time"
)

type AnomalyDataPayload struct {
	// Energy generated in the respective interval in KWH
	Value float64 `json:"value"`

	// Percentage of availability in given interval
	Availability float64 `json:"availability"`

	// Name of the region in which the wind park is located
	Region string `json:"region"`

	// UUID of wind park the current turbine belongs to
	ParkId string `json:"park_id"`

	// Start of time interval when energy was generated in UTC format
	Timestamp string `json:"timestamp"`

	// Timestamp at which event was processed
	ProcessingTimestamp string `json:"processing_timestamp"`
}

func (p *AnomalyDataPayload) Unmarshal(input []byte) error {
	var temp map[string]interface{}

	err := json.Unmarshal(input, &temp)

	if err != nil {
		return err
	}

	required := []string{"park_id", "processing_timestamp", "value", "availability"}
	for _, key := range required {
		if _, ok := temp[key]; !ok {
			return fmt.Errorf("key not found: %s", key)
		}
	}

	allowed := []string{"value", "availability", "region", "park_id", "timestamp", "processing_timestamp"}
	for key := range temp {
		found := false
		for _, property := range allowed {
			if key == property {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("additional property not allowed: %s", key)
		}
	}
	err = json.Unmarshal(input, p)
	if err != nil {
		return err
	}

	if p.Availability < 0 || p.Availability > 1 {
		return fmt.Errorf("invalid value: k=%s v=%s", "availability", fmt.Sprint(p.Availability))
	}

	if _, err := uuid.FromString(p.ParkId); err != nil {
		return fmt.Errorf("property doesn't have uuid format: k=%s v=%s", "park_id", fmt.Sprint(p.ParkId))
	}

	if _, err := time.Parse("2006-01-02 15:04:05", p.Timestamp); err != nil {
		return fmt.Errorf("property doesn't have date-time format: k=%s v=%s", "timestamp", fmt.Sprint(p.Timestamp))
	}

	if _, err := time.Parse("2006-01-02 15:04:05", p.ProcessingTimestamp); err != nil {
		return fmt.Errorf("property doesn't have date-time format: k=%s v=%s", "processing_timestamp", fmt.Sprint(p.ProcessingTimestamp))
	}

	return nil
}

func (p AnomalyDataPayload) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	if p.Availability < 0 || p.Availability > 1 {
		return nil, fmt.Errorf("invalid value: k=%s v=%s", "availability", fmt.Sprint(p.Availability))
	}

	return bytes, nil
}
