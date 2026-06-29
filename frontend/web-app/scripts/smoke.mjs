const baseUrls = {
  auth: process.env.AUTH_URL ?? "http://127.0.0.1:8081",
  wake: process.env.WAKE_URL ?? "http://127.0.0.1:8082",
  device: process.env.DEVICE_URL ?? "http://127.0.0.1:8083",
  projection: process.env.PROJECTION_URL ?? "http://127.0.0.1:8086",
  web: process.env.WEB_URL ?? "http://127.0.0.1:3002"
};

function log(step, ok, detail = "") {
  const prefix = ok ? "[ok]" : "[fail]";
  console.log(`${prefix} ${step}${detail ? ` - ${detail}` : ""}`);
}

async function request(url, options = {}) {
  const response = await fetch(url, {
    headers: { "Content-Type": "application/json", ...(options.headers ?? {}) },
    ...options
  });
  const text = await response.text();
  let body = null;
  try {
    body = text ? JSON.parse(text) : null;
  } catch {
    body = text;
  }
  return { response, body };
}

async function main() {
  const healthChecks = [
    ["auth-service", `${baseUrls.auth}/health`],
    ["wake-session-service", `${baseUrls.wake}/health`],
    ["device-validation-service", `${baseUrls.device}/health`],
    ["projection-service", `${baseUrls.projection}/health`],
    ["web-app", `${baseUrls.web}`]
  ];

  for (const [name, url] of healthChecks) {
    const { response, body } = await request(url);
    if (!response.ok) {
      throw new Error(`${name} health check failed: ${response.status} ${JSON.stringify(body)}`);
    }
    log(`${name} health`, true);
  }

  const login = await request(`${baseUrls.auth}/login`, {
    method: "POST",
    body: JSON.stringify({ email: "marcus@example.com", password: "morning123" })
  });
  if (!login.response.ok || !login.body?.user_id) {
    throw new Error(`login failed: ${login.response.status} ${JSON.stringify(login.body)}`);
  }
  log("login", true, login.body.user_id);

  const device = await request(`${baseUrls.device}/devices/validate`, {
    method: "POST",
    body: JSON.stringify({
      user_agent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/126.0 Safari/537.36",
      viewport_width: 1440,
      viewport_height: 900,
      screen_width: 1440,
      screen_height: 900,
      has_touch: false,
      had_keyboard_event: true,
      had_pointer_event: true,
      correlation_id: "smoke-correlation-1"
    })
  });
  if (!device.response.ok || !device.body?.validation?.allowed || !device.body?.device_proof_id) {
    throw new Error(`device validation failed: ${device.response.status} ${JSON.stringify(device.body)}`);
  }
  log("desktop validation", true, device.body.device_proof_id);

  const checkIn = await request(`${baseUrls.wake}/wake-sessions/check-in`, {
    method: "POST",
    body: JSON.stringify({
      user_id: login.body.user_id,
      target_time: "07:00",
      checked_in_at: new Date("2026-06-28T07:06:00.000Z").toISOString(),
      morning_intent: "Finish Kafka outbox implementation",
      device_proof_id: device.body.device_proof_id,
      correlation_id: "smoke-correlation-1"
    })
  });
  if (!checkIn.response.ok || checkIn.body?.status !== "confirmed") {
    throw new Error(`wake check-in failed: ${checkIn.response.status} ${JSON.stringify(checkIn.body)}`);
  }
  log("wake check-in", true, checkIn.body.wake_session_id);

  const history = await request(`${baseUrls.wake}/wake-sessions/history?user_id=${login.body.user_id}`);
  if (!history.response.ok || !Array.isArray(history.body?.sessions) || history.body.sessions.length === 0) {
    throw new Error(`history fetch failed: ${history.response.status} ${JSON.stringify(history.body)}`);
  }
  log("wake history", true, `${history.body.sessions.length} session(s)`);

  const dashboard = await request(`${baseUrls.projection}/dashboard?user_id=${login.body.user_id}`);
  if (!dashboard.response.ok) {
    throw new Error(`dashboard failed: ${dashboard.response.status} ${JSON.stringify(dashboard.body)}`);
  }
  log("projection dashboard", true);

  console.log("Smoke test passed.");
}

main().catch((error) => {
  console.error(error.message);
  process.exit(1);
});

