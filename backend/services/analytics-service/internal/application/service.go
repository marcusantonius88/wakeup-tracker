package application

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/analytics-service/internal/domain"
	"wakeup-tracker/backend/services/analytics-service/internal/ports"
	"wakeup-tracker/backend/shared/events"
	"wakeup-tracker/backend/shared/observability"
)

type Service struct {
	repository  ports.Repository
	publisher   ports.EventPublisher
	idempotency ports.IdempotencyStore
	metrics     *observability.Metrics
}

func NewService(repository ports.Repository, publisher ports.EventPublisher, idempotency ports.IdempotencyStore, metrics *observability.Metrics) *Service {
	return &Service{repository: repository, publisher: publisher, idempotency: idempotency, metrics: metrics}
}

func (s *Service) Process(ctx context.Context, event events.Envelope) (domain.BehaviorMetrics, error) {
	started := time.Now()
	s.metrics.Inc("events_processed_total")
	defer s.metrics.ObserveDuration("event_processing_duration_seconds", started)

	userID, _ := event.Payload["user_id"].(string)
	if userID == "" {
		userID = event.AggregateID
	}
	if event.EventID != "" {
		reserved, err := s.idempotency.Reserve(ctx, "analytics:event:"+event.EventID, 24*time.Hour)
		if err != nil {
			s.metrics.Inc("errors_total")
			return domain.BehaviorMetrics{}, err
		}
		if !reserved {
			return s.repository.Load(ctx, userID)
		}
	}
	current, err := s.repository.Load(ctx, userID)
	if err != nil {
		s.metrics.Inc("errors_total")
		return domain.BehaviorMetrics{}, err
	}

	next, published := domain.Apply(current, event)
	if err := s.repository.Save(ctx, next); err != nil {
		s.metrics.Inc("errors_total")
		return domain.BehaviorMetrics{}, err
	}
	for _, outgoing := range published {
		if err := s.publisher.Publish(ctx, outgoing); err != nil {
			s.metrics.Inc("errors_total")
			return domain.BehaviorMetrics{}, err
		}
	}

	s.metrics.Add("consistency_score_total", next.ConsistencyScore)
	if next.LastRegressionDetected {
		s.metrics.Inc("wakeup_regression_detected_total")
	}
	return next, nil
}

func (s *Service) Get(ctx context.Context, userID string) (domain.BehaviorMetrics, error) {
	return s.repository.Load(ctx, userID)
}
