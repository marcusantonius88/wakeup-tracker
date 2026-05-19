package memory

import (
	"context"
	"sync"
	"time"

	"wakeup-tracker/backend/services/projection-service/internal/domain"
	"wakeup-tracker/backend/shared/events"
)

type Repository struct {
	mu             sync.Mutex
	wakeSessions   map[string][]domain.WakeSessionProjection
	morningIntents map[string][]domain.MorningIntentProjection
	streaks        map[string]domain.StreakProjection
	timelines      map[string][]events.Envelope
}

func NewRepository() *Repository {
	return &Repository{
		wakeSessions:   map[string][]domain.WakeSessionProjection{},
		morningIntents: map[string][]domain.MorningIntentProjection{},
		streaks:        map[string]domain.StreakProjection{},
		timelines:      map[string][]events.Envelope{},
	}
}

func (r *Repository) Apply(ctx context.Context, event events.Envelope) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userID := domain.EventUserID(event)
	r.timelines[userID] = append(r.timelines[userID], event)

	switch event.EventType {
	case events.WakeSessionStarted:
		target, _ := event.Payload["target_time"].(string)
		r.wakeSessions[userID] = append(r.wakeSessions[userID], domain.WakeSessionProjection{
			WakeSessionID: event.AggregateID,
			UserID:        userID,
			Status:        "started",
			TargetTime:    target,
		})
	case events.WakeUpConfirmed, events.WakeUpFailed:
		status := "confirmed"
		if event.EventType == events.WakeUpFailed {
			status = "failed"
		}
		for index := range r.wakeSessions[userID] {
			if r.wakeSessions[userID][index].WakeSessionID == event.AggregateID {
				r.wakeSessions[userID][index].Status = status
				if checkedIn, ok := event.Payload["checked_in_at"].(time.Time); ok {
					r.wakeSessions[userID][index].CheckedInAt = checkedIn
				}
			}
		}
	case events.MorningIntentSubmitted:
		intent, _ := event.Payload["intent_text"].(string)
		r.morningIntents[userID] = append(r.morningIntents[userID], domain.MorningIntentProjection{
			WakeSessionID: event.AggregateID,
			IntentText:    intent,
			CreatedAt:     event.CreatedAt,
		})
	case events.StreakIncreased:
		streak := intFromAny(event.Payload["streak"])
		r.streaks[userID] = domain.StreakProjection{UserID: userID, Streak: streak}
	case events.StreakBroken:
		r.streaks[userID] = domain.StreakProjection{UserID: userID, Streak: 0}
	}
	return nil
}

func (r *Repository) Dashboard(ctx context.Context, userID string) (domain.Dashboard, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return domain.Dashboard{
		WakeSessions:   append([]domain.WakeSessionProjection{}, r.wakeSessions[userID]...),
		MorningIntents: append([]domain.MorningIntentProjection{}, r.morningIntents[userID]...),
		Streak:         r.streaks[userID],
		Timeline:       append([]events.Envelope{}, r.timelines[userID]...),
	}, nil
}

func intFromAny(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case float64:
		return int(typed)
	default:
		return 0
	}
}
