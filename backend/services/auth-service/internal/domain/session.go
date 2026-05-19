package domain

import (
	"errors"
	"strings"
	"time"

	"wakeup-tracker/backend/shared/events"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type Session struct {
	ID           string
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func Authenticate(email, password string) (*User, *Session, []events.Envelope, error) {
	if !strings.Contains(email, "@") {
		return nil, nil, nil, errors.New("valid email is required")
	}
	if len(password) < 6 {
		return nil, nil, nil, errors.New("password must have at least 6 characters")
	}

	user := &User{ID: stableDemoUser(email), Email: strings.ToLower(strings.TrimSpace(email))}
	session := &Session{
		ID:           events.NewID(),
		UserID:       user.ID,
		AccessToken:  "access-" + events.NewID(),
		RefreshToken: "refresh-" + events.NewID(),
		ExpiresAt:    time.Now().UTC().Add(12 * time.Hour),
		CreatedAt:    time.Now().UTC(),
	}

	published := []events.Envelope{
		events.NewEnvelope(events.UserAuthenticated, user.ID, "", map[string]any{"user_id": user.ID, "email": user.Email}),
		events.NewEnvelope(events.SessionCreated, session.ID, "", map[string]any{"user_id": user.ID, "expires_at": session.ExpiresAt}),
	}
	return user, session, published, nil
}

func stableDemoUser(email string) string {
	if strings.EqualFold(email, "marcus@example.com") {
		return "demo-user"
	}
	return "user-" + strings.ReplaceAll(strings.ToLower(strings.TrimSpace(email)), "@", "-")
}
