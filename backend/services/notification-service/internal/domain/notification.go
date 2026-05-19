package domain

import (
	"errors"
	"strings"
	"time"

	"wakeup-tracker/backend/shared/events"
)

type Notification struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	WakeSessionID string    `json:"wake_session_id"`
	Channel       string    `json:"channel"`
	Message       string    `json:"message"`
	SentAt        time.Time `json:"sent_at"`
}

func Send(userID, wakeSessionID, channel, message, correlationID string) (*Notification, events.Envelope, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, events.Envelope{}, errors.New("user_id is required")
	}
	if channel == "" {
		channel = "browser"
	}
	if message == "" {
		message = "Good morning. Open WakeUpTracker and submit your Morning Intent."
	}
	notification := &Notification{
		ID:            events.NewID(),
		UserID:        userID,
		WakeSessionID: wakeSessionID,
		Channel:       channel,
		Message:       message,
		SentAt:        time.Now().UTC(),
	}
	event := events.NewEnvelope(events.NotificationSent, notification.ID, correlationID, map[string]any{
		"user_id":         notification.UserID,
		"wake_session_id": notification.WakeSessionID,
		"channel":         notification.Channel,
		"sent_at":         notification.SentAt,
	})
	return notification, event, nil
}
