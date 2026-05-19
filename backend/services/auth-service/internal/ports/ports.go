package ports

import (
	"context"

	"wakeup-tracker/backend/services/auth-service/internal/domain"
	"wakeup-tracker/backend/shared/events"
)

type SessionRepository interface {
	Save(ctx context.Context, user *domain.User, session *domain.Session) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event events.Envelope) error
}
