package storage

type AverageProductionEntry struct {
	ParkId string `json:"ParkId"`

	Region string `json:"Region"`

	AverageValue float64 `json:"AverageValue"`

	// Count for AverageValue, used for materialized views
	count_AverageValue int

	// Sum for AverageValue, used for materialized views
	sum_AverageValue float64

	HighestValue float64 `json:"HighestValue"`

	LowestValue float64 `json:"LowestValue"`
}

type AnomalyDetectionEntry struct {
	ParkId string `json:"ParkId"`

	TurbineId string `json:"TurbineId"`

	Value float64 `json:"Value"`

	Availability int `json:"Availability"`
}
