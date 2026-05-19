package domain

import (
	"time"

	"wakeup-tracker/backend/shared/events"
)

type WakeSessionProjection struct {
	WakeSessionID string    `json:"wake_session_id"`
	UserID        string    `json:"user_id"`
	Status        string    `json:"status"`
	TargetTime    string    `json:"target_time"`
	CheckedInAt   time.Time `json:"checked_in_at"`
}

type MorningIntentProjection struct {
	WakeSessionID string    `json:"wake_session_id"`
	IntentText    string    `json:"intent_text"`
	CreatedAt     time.Time `json:"created_at"`
}

type StreakProjection struct {
	UserID string `json:"user_id"`
	Streak int    `json:"streak"`
}

type Dashboard struct {
	WakeSessions   []WakeSessionProjection   `json:"wake_session_projection"`
	MorningIntents []MorningIntentProjection `json:"morning_intent_projection"`
	Streak         StreakProjection          `json:"streak_projection"`
	Timeline       []events.Envelope         `json:"wake_up_timeline_projection"`
}

func EventUserID(event events.Envelope) string {
	if userID, ok := event.Payload["user_id"].(string); ok {
		return userID
	}
	return event.AggregateID
}
