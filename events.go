package go_piggy

import "github.com/twinj/uuid"

type Event struct {
	ID, Type string
	Facts    []Fact
}

// NewEvent creates a new event with a newly generated ID
func NewEvent(eventType string, facts []Fact) Event {
	u := uuid.NewV4()
	return Event{
		ID:    u.String(),
		Type:  eventType,
		Facts: facts,
	}
}

type Fact struct {
	Op, Key, Value string
}

type Emitter interface {
	Send(Event)
}

type Receiver interface {
	Listen() <-chan Event
}
