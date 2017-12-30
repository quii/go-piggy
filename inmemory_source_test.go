package go_piggy

import (
	"testing"
	"reflect"
	"time"
	"fmt"
)

func TestItCanSendAndReceiveFacts(t *testing.T) {

	fact1 := Fact{
		Op:    "SET",
		Key:   "name",
		Value: "cj",
	}

	event1 := NewEvent("manuscript", []Fact{fact1})

	source := NewInMemoryEventSource()

	eventsChannel := source.Listen(0)

	source.Send(event1)

	event, err := waitForEvent(eventsChannel)

	if err != nil {
		t.Fatal(err)
	}

	assertEventEquals(t, event1, event)
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
		return Event{}, fmt.Errorf("Timed out waiting for event")
	}
}
