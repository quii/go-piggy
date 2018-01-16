package go_piggy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestItCanSendAndReceiveFacts(t *testing.T) {

	event := RandomEvent()

	source := NewInMemoryEventSource()

	eventsChannel := source.Listen(0)

	source.Send(event)

	receivedEvent, err := waitForEvent(eventsChannel)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, event, receivedEvent)
}

func TestItCanListenFromAPoint(t *testing.T) {
	event1 := RandomEvent()
	event2 := RandomEvent()
	event3 := RandomEvent()

	source := NewInMemoryEventSource()

	source.Send(event1)
	source.Send(event2)
	source.Send(event3)

	eventsChannel := source.Listen(1)

	firstReceivedEvent, err := waitForEvent(eventsChannel)

	if err != nil {
		t.Fatalf("err waiting for first event, %s", err)
	}

	secondReceivedEvent, err := waitForEvent(eventsChannel)

	if err != nil {
		t.Fatalf("err waiting for second event, %s", err)
	}

	assert.Equal(t, firstReceivedEvent, event2)
	assert.Equal(t, secondReceivedEvent, event3)

}

func waitForEvent(ch <-chan Event) (Event, error) {
	select {
	case e := <-ch:
		return e, nil
	case <-time.After(5 * time.Millisecond):
		return Event{}, fmt.Errorf("timed out waiting for event")
	}
}
