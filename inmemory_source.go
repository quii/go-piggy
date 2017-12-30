package go_piggy

type InMemorySource struct {
	events    []Event
	listeners []chan Event
}

func NewInMemoryEventSource() *InMemorySource {
	return &InMemorySource{}
}

func (i *InMemorySource) Listen(from int) <-chan Event {
	newListener := make(chan Event, 10)

	go i.sendEventsFrom(from, newListener)

	i.listeners = append(i.listeners, newListener)
	return newListener
}

func (i *InMemorySource) sendEventsFrom(from int, to chan<- Event) {
	for _, e := range i.events[from:] {
		to <- e
	}
}

func (i *InMemorySource) Send(event Event) {
	i.events = append(i.events, event)

	//todo: what if the listener's chan is full? fire off a go routine to wait on it? or something else?
	for _, listeners := range i.listeners {
		listeners <- event
	}
}
