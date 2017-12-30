package go_piggy

type InMemorySource struct {
	events    []Event
	listeners [] chan Event
}

func NewInMemoryEventSource() *InMemorySource {
	return &InMemorySource{}
}

func (i *InMemorySource) Listen(from int) <-chan Event {
	newListener := make(chan Event, 10)

	i.listeners = append(i.listeners, newListener)
	return newListener
}

func (i *InMemorySource) Send(event Event) {
	i.events = append(i.events, event)

	for _, listeners := range i.listeners {
		listeners <- event
	}
}
