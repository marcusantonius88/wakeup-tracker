CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS auth_sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  access_token TEXT NOT NULL,
  refresh_token TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS wake_sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  target_time TEXT NOT NULL,
  morning_intent TEXT NOT NULL,
  checked_in_at TIMESTAMPTZ NOT NULL,
  status TEXT NOT NULL,
  streak INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS processed_events (
  event_id TEXT PRIMARY KEY,
  event_type TEXT NOT NULL,
  processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS wake_session_projection (
  wake_session_id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  target_time TEXT NOT NULL,
  status TEXT NOT NULL,
  checked_in_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_consistency_projection (
  user_id TEXT PRIMARY KEY,
  confirmed INTEGER NOT NULL DEFAULT 0,
  failed INTEGER NOT NULL DEFAULT 0,
  consistency_score NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS wake_up_timeline_projection (
  event_id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS streak_projection (
  user_id TEXT PRIMARY KEY,
  streak INTEGER NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS morning_intent_projection (
  wake_session_id TEXT PRIMARY KEY,
  intent_text TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

