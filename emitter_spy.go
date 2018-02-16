package go_piggy

import "reflect"

type EmitterSpy struct {
	events []Event
}

func NewEmitterSpy() *EmitterSpy {
	return &EmitterSpy{}
}

func (i *EmitterSpy) Send(event Event) {
	i.events = append(i.events, event)
}

func (i *EmitterSpy) EventReceived(event Event) bool {
	for _, e := range i.events {
		if reflect.DeepEqual(e, event) {
			return true
		}
	}

	return false
}
