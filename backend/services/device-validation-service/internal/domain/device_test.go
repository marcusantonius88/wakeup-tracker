package domain

import (
	"testing"

	"wakeup-tracker/backend/shared/events"
)

func TestValidateRejectsMobile(t *testing.T) {
	result, event := Validate(ValidationInput{
		UserAgent:        "Mozilla/5.0 iPhone Mobile",
		ViewportWidth:    390,
		ScreenWidth:      390,
		HadPointerEvent:  true,
		HadKeyboardEvent: false,
	})
	if result.Allowed {
		t.Fatal("expected mobile device to be rejected")
	}
	if event.EventType != events.DeviceRejected {
		t.Fatalf("expected DeviceRejected event, got %s", event.EventType)
	}
}

func TestValidateAcceptsDesktopWithInteraction(t *testing.T) {
	result, event := Validate(ValidationInput{
		UserAgent:       "Mozilla/5.0 X11 Linux x86_64",
		ViewportWidth:   1440,
		ViewportHeight:  900,
		ScreenWidth:     1440,
		ScreenHeight:    900,
		HadPointerEvent: true,
	})
	if !result.Allowed {
		t.Fatalf("expected desktop to pass, got reasons: %v", result.Reasons)
	}
	if event.EventType != events.DeviceValidated {
		t.Fatalf("expected DeviceValidated event, got %s", event.EventType)
	}
}
