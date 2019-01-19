package event

import "time"

// Message Key API
type Message interface {
	Key() string
}

// Message API
type MeowCreatedMessage struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

// Prompts response when a meow created.
func (m *MeowCreatedMessage) Key() string {
	return "meow.created"
}
