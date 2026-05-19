export type LoginResponse = {
  user_id: string;
  access_token: string;
  refresh_token: string;
};

export type DeviceValidationResponse = {
  validation: {
    allowed: boolean;
    device_type: string;
    reasons: string[];
    correlation_id: string;
  };
  device_proof_id: string;
};

export type WakeSessionResponse = {
  wake_session_id: string;
  status: string;
  streak: number;
  checked_in_at: string;
  correlation_id: string;
};

const endpoints = {
  auth: process.env.NEXT_PUBLIC_AUTH_API ?? "http://localhost:8081",
  wake: process.env.NEXT_PUBLIC_WAKE_API ?? "http://localhost:8082",
  device: process.env.NEXT_PUBLIC_DEVICE_API ?? "http://localhost:8083",
  projection: process.env.NEXT_PUBLIC_PROJECTION_API ?? "http://localhost:8086"
};

export function newCorrelationId() {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return `corr-${Date.now()}`;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
  const response = await fetch(`${endpoints.auth}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  });
  if (!response.ok) {
    throw new Error("Login failed");
  }
  return response.json();
}

export async function validateDevice(payload: Record<string, unknown>): Promise<DeviceValidationResponse> {
  const response = await fetch(`${endpoints.device}/devices/validate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });
  const data = await response.json();
  if (!response.ok) {
    return data;
  }
  return data;
}

export async function checkIn(payload: Record<string, unknown>): Promise<WakeSessionResponse> {
  const response = await fetch(`${endpoints.wake}/wake-sessions/check-in`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error ?? "Wake check-in failed");
  }
  return data;
}

export async function loadHistory(userId: string) {
  const response = await fetch(`${endpoints.wake}/wake-sessions/history?user_id=${userId}`, {
    cache: "no-store"
  });
  if (!response.ok) {
    return { sessions: [] };
  }
  return response.json();
}

