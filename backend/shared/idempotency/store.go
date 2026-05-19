package idempotency

import (
	"context"
	"sync"
	"time"
)

type Store interface {
	Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error)
}

type MemoryStore struct {
	mu      sync.Mutex
	entries map[string]time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{entries: map[string]time.Time{}}
}

func (s *MemoryStore) Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	for entry, expiresAt := range s.entries {
		if now.After(expiresAt) {
			delete(s.entries, entry)
		}
	}

	if expiresAt, ok := s.entries[key]; ok && now.Before(expiresAt) {
		return false, nil
	}

	s.entries[key] = now.Add(ttl)
	return true, nil
}
