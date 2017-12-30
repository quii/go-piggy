package go_piggy

type Emitter interface {
	Send(Event)
}

type Receiver interface {
	Listen(from int) <-chan Event
}
