package go_piggy

import (
	"testing"
	"reflect"
	"time"
	"fmt"
)

func TestItCanSendAndReceiveFacts(t *testing.T) {

	event := RandomEvent()

	source := NewInMemoryEventSource()

	eventsChannel := source.Listen(0)

	source.Send(event)

	event, err := waitForEvent(eventsChannel)

	if err != nil {
		t.Fatal(err)
	}

	assertEventEquals(t, event, event)
}

func assertEventEquals(t *testing.T, expected, actual Event) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func waitForEvent(ch <-chan Event) (Event, error) {
	select {
	case e := <-ch:
		return e, nil
	case <-time.After(5 * time.Millisecond):
		return Event{}, fmt.Errorf("timed out waiting for event")
	}
}
