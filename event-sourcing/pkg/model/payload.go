package model

type Payload interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}
