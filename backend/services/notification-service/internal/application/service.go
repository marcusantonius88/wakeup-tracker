package application

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/notification-service/internal/domain"
	"wakeup-tracker/backend/services/notification-service/internal/ports"
	"wakeup-tracker/backend/shared/contracts"
	"wakeup-tracker/backend/shared/observability"
)

type Service struct {
	repository ports.NotificationRepository
	publisher  ports.EventPublisher
	metrics    *observability.Metrics
}

func NewService(repository ports.NotificationRepository, publisher ports.EventPublisher, metrics *observability.Metrics) *Service {
	return &Service{repository: repository, publisher: publisher, metrics: metrics}
}

func (s *Service) Send(ctx context.Context, request contracts.NotificationRequest) (*domain.Notification, error) {
	started := time.Now()
	s.metrics.Inc("events_processed_total")
	defer s.metrics.ObserveDuration("event_processing_duration_seconds", started)

	notification, event, err := domain.Send(request.UserID, request.WakeSessionID, request.Channel, request.Message, request.CorrelationID)
	if err != nil {
		s.metrics.Inc("notification_failures_total")
		s.metrics.Inc("errors_total")
		return nil, err
	}
	if err := s.repository.Save(ctx, notification); err != nil {
		s.metrics.Inc("notification_failures_total")
		s.metrics.Inc("errors_total")
		return nil, err
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.metrics.Inc("notification_failures_total")
		s.metrics.Inc("errors_total")
		return nil, err
	}
	s.metrics.Inc("notifications_sent_total")
	return notification, nil
}

func (s *Service) List(ctx context.Context, userID string) ([]domain.Notification, error) {
	return s.repository.List(ctx, userID)
}
