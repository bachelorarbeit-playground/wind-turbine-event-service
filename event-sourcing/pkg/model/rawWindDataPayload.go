package model

import (
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"

	"regexp"

	"time"
)

type RawWindDataPayload struct {
	// Energy generated in the respective interval in MWH
	Value float64 `json:"value"`

	// Percentage of availability in given interval (0-100 integer)
	Availability int `json:"availability"`

	// Name of the region in which the wind park is located
	Region string `json:"region"`

	// UUID of wind park the current turbine belongs to
	ParkId string `json:"park_id"`

	// UUID of the current turbine
	TurbineId string `json:"turbine_id"`

	// Date when energy was generated  (e.g. 11-03-2020)
	Date string `json:"date"`

	// Order of the interval of time in the day (values range from 1-24 for 1 hour intervals)
	Interval int `json:"interval"`

	// Timezone in which the turbine finds itself
	Timezone string `json:"timezone"`
}

func (p *RawWindDataPayload) Unmarshal(input []byte) error {
	var temp map[string]interface{}

	err := json.Unmarshal(input, &temp)

	if err != nil {
		return err
	}

	required := []string{"park_id", "turbine_id", "date", "interval", "value", "availability", "region", "timezone"}
	for _, key := range required {
		if _, ok := temp[key]; !ok {
			return fmt.Errorf("key not found: %s", key)
		}
	}

	allowed := []string{"value", "availability", "region", "park_id", "turbine_id", "date", "interval", "timezone"}
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

	if p.Value < 0 {
		return fmt.Errorf("invalid value: k=%s v=%s", "value", fmt.Sprint(p.Value))
	}

	if p.Availability < 0 || p.Availability > 100 {
		return fmt.Errorf("invalid value: k=%s v=%s", "availability", fmt.Sprint(p.Availability))
	}

	if p.Interval < 1 || p.Interval > 24 {
		return fmt.Errorf("invalid value: k=%s v=%s", "interval", fmt.Sprint(p.Interval))
	}

	if ok, err := regexp.MatchString(`Ber`, p.Region); err != nil {
		return fmt.Errorf("invalid pattern: %s", err)
	} else if !ok {
		return fmt.Errorf("value doesn't match pattern: k=%s v=%s, p=%s", "region", fmt.Sprint(p.Region), `Ber`)
	}

	if len(p.ParkId) < 0 || len(p.ParkId) > 36 {
		return fmt.Errorf("invalid value: k=%s v=%s", "park_id", fmt.Sprint(p.ParkId))
	}

	if _, err := uuid.FromString(p.ParkId); err != nil {
		return fmt.Errorf("property doesn't have uuid format: k=%s v=%s", "park_id", fmt.Sprint(p.ParkId))
	}

	if len(p.TurbineId) < 0 || len(p.TurbineId) > 36 {
		return fmt.Errorf("invalid value: k=%s v=%s", "turbine_id", fmt.Sprint(p.TurbineId))
	}

	if _, err := uuid.FromString(p.TurbineId); err != nil {
		return fmt.Errorf("property doesn't have uuid format: k=%s v=%s", "turbine_id", fmt.Sprint(p.TurbineId))
	}

	if _, err := time.Parse("2006-01-02", p.Date); err != nil {
		return fmt.Errorf("property doesn't have date format: k=%s v=%s", "date", fmt.Sprint(p.Date))
	}

	return nil
}

func (p RawWindDataPayload) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	if p.Value < 0 {
		return nil, fmt.Errorf("invalid value: k=%s v=%s", "value", fmt.Sprint(p.Value))
	}

	if p.Availability < 0 || p.Availability > 100 {
		return nil, fmt.Errorf("invalid value: k=%s v=%s", "availability", fmt.Sprint(p.Availability))
	}

	if p.Interval < 1 || p.Interval > 24 {
		return nil, fmt.Errorf("invalid value: k=%s v=%s", "interval", fmt.Sprint(p.Interval))
	}

	if ok, err := regexp.MatchString(`Ber`, p.Region); err != nil {
		return nil, fmt.Errorf("invalid pattern: %s", err)
	} else if !ok {
		return nil, fmt.Errorf("value doesn't match pattern: k=%s v=%s, p=%s", "region", fmt.Sprint(p.Region), `Ber`)
	}

	if len(p.ParkId) < 0 || len(p.ParkId) > 36 {
		return nil, fmt.Errorf("invalid value: k=%s v=%s", "park_id", fmt.Sprint(p.ParkId))
	}

	if len(p.TurbineId) < 0 || len(p.TurbineId) > 36 {
		return nil, fmt.Errorf("invalid value: k=%s v=%s", "turbine_id", fmt.Sprint(p.TurbineId))
	}

	return bytes, nil
}
