package infra

import (
	"context"
	"encoding/json"
	"sync"

	"wakeup-tracker/backend/shared/events"
	"wakeup-tracker/backend/shared/observability"
)

type EventPublisher interface {
	Publish(ctx context.Context, event events.Envelope) error
}

type InMemoryEventBus struct {
	mu     sync.Mutex
	events []events.Envelope
	logger *observability.Logger
}

func NewInMemoryEventBus(logger *observability.Logger) *InMemoryEventBus {
	return &InMemoryEventBus{logger: logger, events: []events.Envelope{}}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, event events.Envelope) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.events = append(b.events, event)
	payload, _ := json.Marshal(event.Payload)
	b.logger.Info("event published", observability.LogEntry{
		EventID:       event.EventID,
		CorrelationID: event.CorrelationID,
		AggregateID:   event.AggregateID,
		Fields: map[string]any{
			"event_type": event.EventType,
			"payload":    string(payload),
			"transport":  "in-memory-local-dev",
		},
	})
	return nil
}

func (b *InMemoryEventBus) Events() []events.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()
	copyOfEvents := make([]events.Envelope, len(b.events))
	copy(copyOfEvents, b.events)
	return copyOfEvents
}
