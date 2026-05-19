package ports

import (
	"context"

	"wakeup-tracker/backend/shared/events"
)

type EventPublisher interface {
	Publish(ctx context.Context, event events.Envelope) error
}
