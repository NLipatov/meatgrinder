package dtos

type EventDTO struct {
	EventType string      `json:"event_type"`
	Payload   interface{} `json:"payload"`
}
