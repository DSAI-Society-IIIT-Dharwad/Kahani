package observer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// EventType enumerates notifications emitted by the system.
type EventType string

const (
	EventBlockCommitted       EventType = "block.committed"
	EventTransactionQueued    EventType = "transaction.queued"
	EventTransactionCommitted EventType = "transaction.committed"
	EventError                EventType = "error"
)

// Event represents a message broadcast to subscribers.
type Event struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Bus fan-outs events to registered subscribers.
type Bus struct {
	mu      sync.RWMutex
	subs    map[string]chan Event
	counter uint64
	closed  bool
}

// NewBus constructs an empty event bus.
func NewBus() *Bus {
	return &Bus{subs: make(map[string]chan Event)}
}

// Subscribe registers a new subscriber and returns its identifier and channel.
// A default buffer of 16 is applied when buffer <= 0.
func (b *Bus) Subscribe(buffer int) (string, <-chan Event) {
	if b == nil {
		ch := make(chan Event)
		close(ch)
		return "", ch
	}

	if buffer <= 0 {
		buffer = 16
	}

	ch := make(chan Event, buffer)

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		close(ch)
		return "", ch
	}

	id := fmt.Sprintf("sub-%d", atomic.AddUint64(&b.counter, 1))
	b.subs[id] = ch
	return id, ch
}

// Unsubscribe removes the subscriber and closes its channel.
func (b *Bus) Unsubscribe(id string) {
	if b == nil || id == "" {
		return
	}

	var ch chan Event

	b.mu.Lock()
	if b.subs != nil {
		ch = b.subs[id]
		delete(b.subs, id)
	}
	b.mu.Unlock()

	if ch != nil {
		close(ch)
	}
}

// Publish fan-outs the event to all subscribers using best-effort delivery.
func (b *Bus) Publish(event Event) {
	if b == nil {
		return
	}

	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return
	}

	targets := make([]chan Event, 0, len(b.subs))
	for _, ch := range b.subs {
		targets = append(targets, ch)
	}
	b.mu.RUnlock()

	for _, ch := range targets {
		select {
		case ch <- event:
		default:
		}
	}
}

// Close shuts down the bus and closes all subscriber channels.
func (b *Bus) Close() {
	if b == nil {
		return
	}

	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.closed = true

	for id, ch := range b.subs {
		close(ch)
		delete(b.subs, id)
	}

	b.mu.Unlock()
}
