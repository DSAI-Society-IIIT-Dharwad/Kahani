package observer

import (
	"sync"
	"testing"
	"time"
)

func TestBusPublishAndSubscribe(t *testing.T) {
	bus := NewBus()
	id, ch := bus.Subscribe(0)
	if id == "" {
		t.Fatalf("expected subscriber id to be returned")
	}

	event := Event{Type: EventBlockCommitted, Timestamp: time.Now()}
	bus.Publish(event)

	select {
	case received := <-ch:
		if received.Type != event.Type {
			t.Fatalf("expected event type %s, got %s", event.Type, received.Type)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for event")
	}

	bus.Unsubscribe(id)
}

func TestBusUnsubscribeClosesChannel(t *testing.T) {
	bus := NewBus()
	id, ch := bus.Subscribe(1)
	bus.Unsubscribe(id)

	if _, ok := <-ch; ok {
		t.Fatalf("expected channel to be closed")
	}
}

func TestBusClose(t *testing.T) {
	bus := NewBus()
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		_, ch := bus.Subscribe(1)
		wg.Add(1)
		go func(c <-chan Event) {
			defer wg.Done()
			<-c
		}(ch)
	}

	bus.Close()
	wg.Wait()

	// Publish after close should be a no-op.
	bus.Publish(Event{Type: EventError})
}
