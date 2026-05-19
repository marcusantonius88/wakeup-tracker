package contracts

import "time"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DeviceValidationRequest struct {
	UserAgent        string `json:"user_agent"`
	ViewportWidth    int    `json:"viewport_width"`
	ViewportHeight   int    `json:"viewport_height"`
	ScreenWidth      int    `json:"screen_width"`
	ScreenHeight     int    `json:"screen_height"`
	HasTouch         bool   `json:"has_touch"`
	HadKeyboardEvent bool   `json:"had_keyboard_event"`
	HadPointerEvent  bool   `json:"had_pointer_event"`
	CorrelationID    string `json:"correlation_id"`
}

type DeviceValidationResponse struct {
	Allowed       bool     `json:"allowed"`
	DeviceType    string   `json:"device_type"`
	Reasons       []string `json:"reasons"`
	CorrelationID string   `json:"correlation_id"`
}

type WakeCheckInRequest struct {
	UserID        string    `json:"user_id"`
	TargetTime    string    `json:"target_time"`
	CheckedInAt   time.Time `json:"checked_in_at"`
	MorningIntent string    `json:"morning_intent"`
	DeviceProofID string    `json:"device_proof_id"`
	CorrelationID string    `json:"correlation_id"`
}

type WakeSessionResponse struct {
	WakeSessionID string    `json:"wake_session_id"`
	Status        string    `json:"status"`
	Streak        int       `json:"streak"`
	CheckedInAt   time.Time `json:"checked_in_at"`
	CorrelationID string    `json:"correlation_id"`
}

type NotificationRequest struct {
	UserID        string `json:"user_id"`
	WakeSessionID string `json:"wake_session_id"`
	Channel       string `json:"channel"`
	Message       string `json:"message"`
	CorrelationID string `json:"correlation_id"`
}
