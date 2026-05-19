package memory

import (
	"context"
	"sync"

	"wakeup-tracker/backend/services/auth-service/internal/domain"
)

type SessionRepository struct {
	mu       sync.Mutex
	users    map[string]domain.User
	sessions map[string]domain.Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{users: map[string]domain.User{}, sessions: map[string]domain.Session{}}
}

func (r *SessionRepository) Save(ctx context.Context, user *domain.User, session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = *user
	r.sessions[session.ID] = *session
	return nil
}
