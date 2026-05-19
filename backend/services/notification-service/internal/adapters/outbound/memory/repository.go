package memory

import (
	"context"
	"sync"

	"wakeup-tracker/backend/services/notification-service/internal/domain"
)

type NotificationRepository struct {
	mu            sync.Mutex
	notifications []domain.Notification
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{notifications: []domain.Notification{}}
}

func (r *NotificationRepository) Save(ctx context.Context, notification *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications = append(r.notifications, *notification)
	return nil
}

func (r *NotificationRepository) List(ctx context.Context, userID string) ([]domain.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := []domain.Notification{}
	for _, notification := range r.notifications {
		if notification.UserID == userID {
			result = append(result, notification)
		}
	}
	return result, nil
}
