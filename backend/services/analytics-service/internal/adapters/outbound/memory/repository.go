package memory

import (
	"context"
	"sync"

	"wakeup-tracker/backend/services/analytics-service/internal/domain"
)

type Repository struct {
	mu      sync.Mutex
	metrics map[string]domain.BehaviorMetrics
}

func NewRepository() *Repository {
	return &Repository{metrics: map[string]domain.BehaviorMetrics{}}
}

func (r *Repository) Load(ctx context.Context, userID string) (domain.BehaviorMetrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if value, ok := r.metrics[userID]; ok {
		return value, nil
	}
	return domain.BehaviorMetrics{UserID: userID}, nil
}

func (r *Repository) Save(ctx context.Context, metrics domain.BehaviorMetrics) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[metrics.UserID] = metrics
	return nil
}
