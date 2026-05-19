# WakeUpTracker Kafka Topics

- `auth.events`
- `wake-session.events`
- `device-validation.events`
- `notification.events`
- `analytics.events`
- `projection.replay`

The Go services expose hexagonal event publisher ports. The local adapter stores events in memory so the project can run without broker credentials during early development; Kafka wiring belongs behind the same outbound adapter boundary.

