package domain

import (
	"errors"
	"strings"
	"time"

	"wakeup-tracker/backend/shared/events"
)

const (
	StatusStarted   = "started"
	StatusConfirmed = "confirmed"
	StatusFailed    = "failed"
)

type WakeSession struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	TargetTime    string    `json:"target_time"`
	MorningIntent string    `json:"morning_intent"`
	CheckedInAt   time.Time `json:"checked_in_at"`
	Status        string    `json:"status"`
	Streak        int       `json:"streak"`
	CreatedAt     time.Time `json:"created_at"`
}

type CheckInCommand struct {
	UserID        string
	TargetTime    string
	CheckedInAt   time.Time
	MorningIntent string
	DeviceProofID string
	CorrelationID string
	CurrentStreak int
}

func NewWakeSession(cmd CheckInCommand) (*WakeSession, []events.Envelope, error) {
	if strings.TrimSpace(cmd.UserID) == "" {
		return nil, nil, errors.New("user_id is required")
	}
	if strings.TrimSpace(cmd.DeviceProofID) == "" {
		return nil, nil, errors.New("desktop device proof is required before wake check-in")
	}
	if strings.TrimSpace(cmd.MorningIntent) == "" {
		return nil, nil, errors.New("morning intent is required before wake-up confirmation")
	}
	if len(strings.TrimSpace(cmd.MorningIntent)) < 8 {
		return nil, nil, errors.New("morning intent must contain a meaningful objective")
	}

	checkedInAt := cmd.CheckedInAt
	if checkedInAt.IsZero() {
		checkedInAt = time.Now().UTC()
	}

	id := events.NewID()
	session := &WakeSession{
		ID:            id,
		UserID:        cmd.UserID,
		TargetTime:    fallbackTarget(cmd.TargetTime),
		MorningIntent: strings.TrimSpace(cmd.MorningIntent),
		CheckedInAt:   checkedInAt.UTC(),
		Status:        StatusStarted,
		CreatedAt:     time.Now().UTC(),
	}

	session.Status = StatusConfirmed
	session.Streak = cmd.CurrentStreak + 1
	if missedTolerance(session.TargetTime, session.CheckedInAt, 10*time.Minute) {
		session.Status = StatusFailed
		session.Streak = 0
	}

	published := []events.Envelope{
		events.NewEnvelope(events.WakeSessionStarted, session.ID, cmd.CorrelationID, map[string]any{
			"user_id":     session.UserID,
			"target_time": session.TargetTime,
		}),
		events.NewEnvelope(events.MorningIntentSubmitted, session.ID, cmd.CorrelationID, map[string]any{
			"wake_session_id": session.ID,
			"intent_text":     session.MorningIntent,
			"created_at":      session.CreatedAt,
		}),
	}

	if session.Status == StatusConfirmed {
		published = append(published,
			events.NewEnvelope(events.WakeUpConfirmed, session.ID, cmd.CorrelationID, map[string]any{
				"user_id":       session.UserID,
				"checked_in_at": session.CheckedInAt,
				"streak":        session.Streak,
			}),
			events.NewEnvelope(events.StreakIncreased, session.UserID, cmd.CorrelationID, map[string]any{
				"user_id": session.UserID,
				"streak":  session.Streak,
			}),
		)
		return session, published, nil
	}

	published = append(published,
		events.NewEnvelope(events.WakeUpFailed, session.ID, cmd.CorrelationID, map[string]any{
			"user_id":       session.UserID,
			"checked_in_at": session.CheckedInAt,
			"target_time":   session.TargetTime,
			"reason":        "outside_tolerance_window",
		}),
		events.NewEnvelope(events.StreakBroken, session.UserID, cmd.CorrelationID, map[string]any{
			"user_id": session.UserID,
		}),
	)
	return session, published, nil
}

func fallbackTarget(target string) string {
	if strings.TrimSpace(target) == "" {
		return "07:00"
	}
	return strings.TrimSpace(target)
}

func missedTolerance(target string, checkedInAt time.Time, tolerance time.Duration) bool {
	parsed, err := time.Parse("15:04", target)
	if err != nil {
		return false
	}

	deadline := time.Date(
		checkedInAt.Year(),
		checkedInAt.Month(),
		checkedInAt.Day(),
		parsed.Hour(),
		parsed.Minute(),
		0,
		0,
		checkedInAt.Location(),
	).Add(tolerance)

	return checkedInAt.After(deadline)
}
