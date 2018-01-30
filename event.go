package go_piggy

import "github.com/twinj/uuid"

type Event struct {
	Name           string
	EntityID, Type string
	Facts          []Fact
}

type Fact struct {
	Op, Key, Value string
}

func RandomID() string {
	return uuid.NewV4().String()
}

func RandomEvent() Event {
	return Event{
		Name:     "name-" + RandomID(),
		EntityID: "id-" + RandomID(),
		Type:     "random" + RandomID(),
		Facts:    nil,
	}
}
