package application

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/auth-service/internal/domain"
	"wakeup-tracker/backend/services/auth-service/internal/ports"
	"wakeup-tracker/backend/shared/contracts"
	"wakeup-tracker/backend/shared/observability"
)

type Service struct {
	repository ports.SessionRepository
	publisher  ports.EventPublisher
	metrics    *observability.Metrics
}

func NewService(repository ports.SessionRepository, publisher ports.EventPublisher, metrics *observability.Metrics) *Service {
	return &Service{repository: repository, publisher: publisher, metrics: metrics}
}

func (s *Service) Login(ctx context.Context, request contracts.LoginRequest) (contracts.LoginResponse, error) {
	started := time.Now()
	s.metrics.Inc("events_processed_total")
	defer s.metrics.ObserveDuration("event_processing_duration_seconds", started)

	user, session, publishedEvents, err := domain.Authenticate(request.Email, request.Password)
	if err != nil {
		s.metrics.Inc("errors_total")
		return contracts.LoginResponse{}, err
	}

	if err := s.repository.Save(ctx, user, session); err != nil {
		s.metrics.Inc("errors_total")
		return contracts.LoginResponse{}, err
	}

	for _, event := range publishedEvents {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.metrics.Inc("errors_total")
			return contracts.LoginResponse{}, err
		}
	}

	return contracts.LoginResponse{UserID: user.ID, AccessToken: session.AccessToken, RefreshToken: session.RefreshToken}, nil
}
