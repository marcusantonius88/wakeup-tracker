package events

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const (
	UserAuthenticated        = "UserAuthenticated"
	SessionCreated           = "SessionCreated"
	SessionExpired           = "SessionExpired"
	WakeSessionStarted       = "WakeSessionStarted"
	WakeUpConfirmed          = "WakeUpConfirmed"
	WakeUpFailed             = "WakeUpFailed"
	StreakIncreased          = "StreakIncreased"
	StreakBroken             = "StreakBroken"
	MorningIntentSubmitted   = "MorningIntentSubmitted"
	NotificationSent         = "NotificationSent"
	DeviceValidated          = "DeviceValidated"
	DeviceRejected           = "DeviceRejected"
	ConsistencyImproved      = "ConsistencyImproved"
	WakeUpRegressionDetected = "WakeUpRegressionDetected"
)

type Envelope struct {
	EventID       string         `json:"event_id"`
	EventType     string         `json:"event_type"`
	AggregateID   string         `json:"aggregate_id"`
	CorrelationID string         `json:"correlation_id"`
	Payload       map[string]any `json:"payload"`
	CreatedAt     time.Time      `json:"created_at"`
}

func NewEnvelope(eventType, aggregateID, correlationID string, payload map[string]any) Envelope {
	if correlationID == "" {
		correlationID = NewID()
	}

	return Envelope{
		EventID:       NewID(),
		EventType:     eventType,
		AggregateID:   aggregateID,
		CorrelationID: correlationID,
		Payload:       payload,
		CreatedAt:     time.Now().UTC(),
	}
}

func NewID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	dst := make([]byte, 36)
	hex.Encode(dst[0:8], bytes[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], bytes[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], bytes[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], bytes[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], bytes[10:])
	return string(dst)
}
