package application

import (
	"context"
	"time"

	"wakeup-tracker/backend/services/device-validation-service/internal/domain"
	"wakeup-tracker/backend/services/device-validation-service/internal/ports"
	"wakeup-tracker/backend/shared/contracts"
	"wakeup-tracker/backend/shared/observability"
)

type Service struct {
	publisher ports.EventPublisher
	metrics   *observability.Metrics
	logger    *observability.Logger
}

func NewService(publisher ports.EventPublisher, metrics *observability.Metrics, logger *observability.Logger) *Service {
	return &Service{publisher: publisher, metrics: metrics, logger: logger}
}

func (s *Service) Validate(ctx context.Context, request contracts.DeviceValidationRequest) (contracts.DeviceValidationResponse, string, error) {
	started := time.Now()
	s.metrics.Inc("events_processed_total")
	defer s.metrics.ObserveDuration("event_processing_duration_seconds", started)

	result, event := domain.Validate(domain.ValidationInput{
		UserAgent:        request.UserAgent,
		ViewportWidth:    request.ViewportWidth,
		ViewportHeight:   request.ViewportHeight,
		ScreenWidth:      request.ScreenWidth,
		ScreenHeight:     request.ScreenHeight,
		HasTouch:         request.HasTouch,
		HadKeyboardEvent: request.HadKeyboardEvent,
		HadPointerEvent:  request.HadPointerEvent,
		CorrelationID:    request.CorrelationID,
	})

	if err := s.publisher.Publish(ctx, event); err != nil {
		s.metrics.Inc("errors_total")
		return contracts.DeviceValidationResponse{}, "", err
	}

	if result.Allowed {
		s.metrics.Inc("desktop_validation_success_total")
	} else {
		s.metrics.Inc("device_rejected_total")
	}

	s.logger.Info("device validation completed", observability.LogEntry{
		EventID:       event.EventID,
		CorrelationID: event.CorrelationID,
		AggregateID:   event.AggregateID,
		Fields: map[string]any{
			"allowed": result.Allowed,
			"reasons": result.Reasons,
		},
	})

	return contracts.DeviceValidationResponse{
		Allowed:       result.Allowed,
		DeviceType:    result.DeviceType,
		Reasons:       result.Reasons,
		CorrelationID: request.CorrelationID,
	}, result.ProofID, nil
}
