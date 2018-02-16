package go_piggy

type Emitter interface {
	Send(Event)
}

type Receiver interface {
	Listen(from int) <-chan Event
}

type EventSource interface {
	Emitter
	Receiver
}

type Aggregate interface {
	ProcessCommand(event Event) (accepted bool)
}
