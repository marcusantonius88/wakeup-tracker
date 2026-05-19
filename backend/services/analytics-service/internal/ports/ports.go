package ports

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/analytics-service/internal/domain"
	"wakeup-tracker/backend/shared/events"
)

type Repository interface {
	Load(ctx context.Context, userID string) (domain.BehaviorMetrics, error)
	Save(ctx context.Context, metrics domain.BehaviorMetrics) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event events.Envelope) error
}

type IdempotencyStore interface {
	Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error)
}
