package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"wakeup-tracker/backend/services/wake-session-service/internal/domain"
	"wakeup-tracker/backend/services/wake-session-service/internal/ports"
	"wakeup-tracker/backend/shared/contracts"
	"wakeup-tracker/backend/shared/observability"
)

type Service struct {
	repository  ports.WakeSessionRepository
	publisher   ports.EventPublisher
	idempotency ports.IdempotencyStore
	metrics     *observability.Metrics
	logger      *observability.Logger
}

func NewService(repository ports.WakeSessionRepository, publisher ports.EventPublisher, idempotency ports.IdempotencyStore, metrics *observability.Metrics, logger *observability.Logger) *Service {
	return &Service{repository: repository, publisher: publisher, idempotency: idempotency, metrics: metrics, logger: logger}
}

func (s *Service) CheckIn(ctx context.Context, request contracts.WakeCheckInRequest) (contracts.WakeSessionResponse, error) {
	started := time.Now()
	s.metrics.Inc("events_processed_total")
	defer s.metrics.ObserveDuration("event_processing_duration_seconds", started)

	key := fmt.Sprintf("wake-check-in:%s:%s", request.UserID, time.Now().UTC().Format("2006-01-02"))
	reserved, err := s.idempotency.Reserve(ctx, key, 24*time.Hour)
	if err != nil {
		s.metrics.Inc("errors_total")
		return contracts.WakeSessionResponse{}, err
	}
	if !reserved {
		s.metrics.Inc("errors_total")
		return contracts.WakeSessionResponse{}, errors.New("wake-up already checked in today")
	}

	currentStreak, err := s.repository.CurrentStreak(ctx, request.UserID)
	if err != nil {
		s.metrics.Inc("errors_total")
		return contracts.WakeSessionResponse{}, err
	}

	session, publishedEvents, err := domain.NewWakeSession(domain.CheckInCommand{
		UserID:        request.UserID,
		TargetTime:    request.TargetTime,
		CheckedInAt:   request.CheckedInAt,
		MorningIntent: request.MorningIntent,
		DeviceProofID: request.DeviceProofID,
		CorrelationID: request.CorrelationID,
		CurrentStreak: currentStreak,
	})
	if err != nil {
		s.metrics.Inc("errors_total")
		return contracts.WakeSessionResponse{}, err
	}

	if err := s.repository.Save(ctx, session); err != nil {
		s.metrics.Inc("errors_total")
		return contracts.WakeSessionResponse{}, err
	}

	for _, event := range publishedEvents {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.metrics.Inc("errors_total")
			return contracts.WakeSessionResponse{}, err
		}
	}

	s.metrics.Inc("morning_intent_submitted_total")
	if session.Status == domain.StatusConfirmed {
		s.metrics.Inc("wakeup_confirmed_total")
		s.metrics.Inc("streak_increased_total")
	} else {
		s.metrics.Inc("wakeup_failed_total")
	}

	s.logger.Info("wake session check-in completed", observability.LogEntry{
		CorrelationID: request.CorrelationID,
		AggregateID:   session.ID,
		Fields: map[string]any{
			"user_id": session.UserID,
			"status":  session.Status,
			"streak":  session.Streak,
		},
	})

	return contracts.WakeSessionResponse{
		WakeSessionID: session.ID,
		Status:        session.Status,
		Streak:        session.Streak,
		CheckedInAt:   session.CheckedInAt,
		CorrelationID: request.CorrelationID,
	}, nil
}

func (s *Service) History(ctx context.Context, userID string) ([]domain.WakeSession, error) {
	return s.repository.ListByUser(ctx, userID)
}
