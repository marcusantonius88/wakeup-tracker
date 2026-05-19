package memory

import (
	"context"
	"sync"

	"wakeup-tracker/backend/services/wake-session-service/internal/domain"
)

type WakeSessionRepository struct {
	mu       sync.Mutex
	sessions []domain.WakeSession
}

func NewWakeSessionRepository() *WakeSessionRepository {
	return &WakeSessionRepository{sessions: []domain.WakeSession{}}
}

func (r *WakeSessionRepository) Save(ctx context.Context, session *domain.WakeSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions = append(r.sessions, *session)
	return nil
}

func (r *WakeSessionRepository) CurrentStreak(ctx context.Context, userID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	streak := 0
	for _, session := range r.sessions {
		if session.UserID == userID && session.Status == domain.StatusConfirmed {
			streak = session.Streak
		}
		if session.UserID == userID && session.Status == domain.StatusFailed {
			streak = 0
		}
	}
	return streak, nil
}

func (r *WakeSessionRepository) ListByUser(ctx context.Context, userID string) ([]domain.WakeSession, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := []domain.WakeSession{}
	for _, session := range r.sessions {
		if session.UserID == userID {
			result = append(result, session)
		}
	}
	return result, nil
}
