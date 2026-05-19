package ports

import (
	"context"

	"wakeup-tracker/backend/services/notification-service/internal/domain"
	"wakeup-tracker/backend/shared/events"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification *domain.Notification) error
	List(ctx context.Context, userID string) ([]domain.Notification, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event events.Envelope) error
}
