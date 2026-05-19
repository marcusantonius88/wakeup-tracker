package application

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/projection-service/internal/domain"
	"wakeup-tracker/backend/services/projection-service/internal/ports"
	"wakeup-tracker/backend/shared/events"
	"wakeup-tracker/backend/shared/observability"
)

type Service struct {
	repository  ports.ProjectionRepository
	idempotency ports.IdempotencyStore
	metrics     *observability.Metrics
}

func NewService(repository ports.ProjectionRepository, idempotency ports.IdempotencyStore, metrics *observability.Metrics) *Service {
	return &Service{repository: repository, idempotency: idempotency, metrics: metrics}
}

func (s *Service) Process(ctx context.Context, event events.Envelope) error {
	started := time.Now()
	s.metrics.Inc("events_processed_total")
	defer s.metrics.ObserveDuration("event_processing_duration_seconds", started)
	if event.EventID != "" {
		reserved, err := s.idempotency.Reserve(ctx, "projection:event:"+event.EventID, 24*time.Hour)
		if err != nil {
			s.metrics.Inc("errors_total")
			return err
		}
		if !reserved {
			return nil
		}
	}
	if err := s.repository.Apply(ctx, event); err != nil {
		s.metrics.Inc("errors_total")
		return err
	}
	return nil
}

func (s *Service) Dashboard(ctx context.Context, userID string) (domain.Dashboard, error) {
	return s.repository.Dashboard(ctx, userID)
}
