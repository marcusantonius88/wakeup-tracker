package ports

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/wake-session-service/internal/domain"
	"wakeup-tracker/backend/shared/events"
)

type WakeSessionRepository interface {
	Save(ctx context.Context, session *domain.WakeSession) error
	CurrentStreak(ctx context.Context, userID string) (int, error)
	ListByUser(ctx context.Context, userID string) ([]domain.WakeSession, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event events.Envelope) error
}

type IdempotencyStore interface {
	Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error)
}
