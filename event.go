package go_piggy

import "github.com/twinj/uuid"

type Event struct {
	ID, Type string
	Facts    []Fact
}

type Fact struct {
	Op, Key, Value string
}

func RandomID() string {
	return uuid.NewV4().String()
}

func RandomEvent() Event {
	return Event{
		ID:    RandomID(),
		Type:  "random"+RandomID(),
		Facts: nil,
	}
}
