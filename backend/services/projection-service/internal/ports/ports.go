package ports

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/projection-service/internal/domain"
	"wakeup-tracker/backend/shared/events"
)

type ProjectionRepository interface {
	Apply(ctx context.Context, event events.Envelope) error
	Dashboard(ctx context.Context, userID string) (domain.Dashboard, error)
}

type IdempotencyStore interface {
	Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error)
}
