package infra

import (
	"context"
	"time"

	"wakeup-tracker/backend/shared/events"
	"wakeup-tracker/backend/shared/observability"
)

type RetryingPublisher struct {
	delegate EventPublisher
	attempts int
	delay    time.Duration
	metrics  *observability.Metrics
	logger   *observability.Logger
}

func NewRetryingPublisher(delegate EventPublisher, attempts int, delay time.Duration, metrics *observability.Metrics, logger *observability.Logger) *RetryingPublisher {
	if attempts < 1 {
		attempts = 1
	}
	return &RetryingPublisher{delegate: delegate, attempts: attempts, delay: delay, metrics: metrics, logger: logger}
}

func (p *RetryingPublisher) Publish(ctx context.Context, event events.Envelope) error {
	var lastErr error
	for attempt := 1; attempt <= p.attempts; attempt++ {
		if err := p.delegate.Publish(ctx, event); err != nil {
			lastErr = err
			p.metrics.Inc("errors_total")
			p.logger.Error("event publish attempt failed", observability.LogEntry{
				EventID:       event.EventID,
				CorrelationID: event.CorrelationID,
				AggregateID:   event.AggregateID,
				Fields: map[string]any{
					"event_type": event.EventType,
					"attempt":    attempt,
				},
			})
			if attempt < p.attempts {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(p.delay):
				}
			}
			continue
		}
		return nil
	}
	return lastErr
}
