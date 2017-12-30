package go_piggy

import "github.com/twinj/uuid"

type Event struct {
	ID, Type string
	Facts    []Fact
}

type Fact struct {
	Op, Key, Value string
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

func RandomEvent() Event {
	u := uuid.NewV4()
	return NewEvent("type"+u.String(), nil)
}
