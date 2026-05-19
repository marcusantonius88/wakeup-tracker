# Redis Usage

Redis is reserved for:

- wake check-in idempotency keys
- event consumer deduplication
- auth refresh/session cache
- retry locks and replay coordination

The scaffold includes an in-memory idempotency implementation with the same port used by the wake-session application layer.

