package domain

import (
	"testing"
	"time"

	"wakeup-tracker/backend/shared/events"
)

func TestNewWakeSessionRequiresMorningIntent(t *testing.T) {
	_, _, err := NewWakeSession(CheckInCommand{
		UserID:        "demo-user",
		TargetTime:    "07:00",
		CheckedInAt:   time.Date(2026, 5, 19, 7, 3, 0, 0, time.UTC),
		DeviceProofID: "desktop-proof",
	})
	if err == nil {
		t.Fatal("expected morning intent validation error")
	}
}

func TestNewWakeSessionRequiresDesktopProof(t *testing.T) {
	_, _, err := NewWakeSession(CheckInCommand{
		UserID:        "demo-user",
		TargetTime:    "07:00",
		CheckedInAt:   time.Date(2026, 5, 19, 7, 3, 0, 0, time.UTC),
		MorningIntent: "Finish Kafka outbox implementation",
	})
	if err == nil {
		t.Fatal("expected desktop proof validation error")
	}
}

func TestNewWakeSessionPublishesMorningIntentBeforeConfirmation(t *testing.T) {
	session, published, err := NewWakeSession(CheckInCommand{
		UserID:        "demo-user",
		TargetTime:    "07:00",
		CheckedInAt:   time.Date(2026, 5, 19, 7, 6, 0, 0, time.UTC),
		MorningIntent: "Finish Kafka outbox implementation",
		DeviceProofID: "desktop-proof",
		CurrentStreak: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if session.Status != StatusConfirmed {
		t.Fatalf("expected confirmed session, got %s", session.Status)
	}
	if published[0].EventType != events.WakeSessionStarted {
		t.Fatalf("expected first event WakeSessionStarted, got %s", published[0].EventType)
	}
	if published[1].EventType != events.MorningIntentSubmitted {
		t.Fatalf("expected second event MorningIntentSubmitted, got %s", published[1].EventType)
	}
	if published[2].EventType != events.WakeUpConfirmed {
		t.Fatalf("expected third event WakeUpConfirmed, got %s", published[2].EventType)
	}
}

func TestNewWakeSessionFailsOutsideTolerance(t *testing.T) {
	session, published, err := NewWakeSession(CheckInCommand{
		UserID:        "demo-user",
		TargetTime:    "07:00",
		CheckedInAt:   time.Date(2026, 5, 19, 7, 20, 0, 0, time.UTC),
		MorningIntent: "Finish Kafka outbox implementation",
		DeviceProofID: "desktop-proof",
		CurrentStreak: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if session.Status != StatusFailed {
		t.Fatalf("expected failed session, got %s", session.Status)
	}
	if published[2].EventType != events.WakeUpFailed {
		t.Fatalf("expected WakeUpFailed event, got %s", published[2].EventType)
	}
}
