package event

type EventStatus int

const (
	StatusPassed EventStatus = iota
	StatusDropped
	StatusFailed
)

type State struct {
	Id       string `json:"id"`
	Instance string `json:"instance"`
	State    string `json:"state"`
	Error    string `json:"error"`
}
