package domain

import (
	"strings"

	"wakeup-tracker/backend/shared/events"
)

type ValidationInput struct {
	UserAgent        string
	ViewportWidth    int
	ViewportHeight   int
	ScreenWidth      int
	ScreenHeight     int
	HasTouch         bool
	HadKeyboardEvent bool
	HadPointerEvent  bool
	CorrelationID    string
}

type ValidationResult struct {
	Allowed    bool
	DeviceType string
	Reasons    []string
	ProofID    string
}

func Validate(input ValidationInput) (ValidationResult, events.Envelope) {
	result := ValidationResult{Allowed: true, DeviceType: "desktop", Reasons: []string{}, ProofID: events.NewID()}
	ua := strings.ToLower(input.UserAgent)

	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		result.Allowed = false
		result.DeviceType = "mobile"
		result.Reasons = append(result.Reasons, "mobile user-agent is blocked in the MVP")
	}
	if strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet") {
		result.Allowed = false
		result.DeviceType = "tablet"
		result.Reasons = append(result.Reasons, "tablet user-agent is blocked in the MVP")
	}
	if input.HasTouch && input.ViewportWidth < 1100 {
		result.Allowed = false
		result.Reasons = append(result.Reasons, "touch-first small viewport suggests mobile/tablet")
	}
	if input.ViewportWidth < 1024 || input.ScreenWidth < 1024 {
		result.Allowed = false
		result.Reasons = append(result.Reasons, "minimum desktop width is 1024px")
	}
	if !input.HadKeyboardEvent && !input.HadPointerEvent {
		result.Allowed = false
		result.Reasons = append(result.Reasons, "keyboard or mouse interaction proof is required")
	}

	eventType := events.DeviceValidated
	if !result.Allowed {
		eventType = events.DeviceRejected
		result.ProofID = ""
	}

	event := events.NewEnvelope(eventType, result.ProofID, input.CorrelationID, map[string]any{
		"allowed":     result.Allowed,
		"device_type": result.DeviceType,
		"reasons":     result.Reasons,
	})
	return result, event
}
