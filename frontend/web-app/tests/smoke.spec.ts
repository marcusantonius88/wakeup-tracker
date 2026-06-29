import { expect, test } from "@playwright/test";

test("wakeup tracker smoke flow", async ({ page }) => {
  const loginRequests: string[] = [];
  const validationRequests: any[] = [];
  const checkInRequests: any[] = [];

  await page.route("**/login", async (route) => {
    loginRequests.push(route.request().postData() ?? "");
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        user_id: "demo-user",
        access_token: "token-access",
        refresh_token: "token-refresh"
      })
    });
  });

  await page.route("**/devices/validate", async (route) => {
    const body = route.request().postDataJSON();
    validationRequests.push(body);
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        validation: {
          allowed: true,
          device_type: "desktop",
          reasons: [],
          correlation_id: body.correlation_id
        },
        device_proof_id: "desktop-proof-123"
      })
    });
  });

  await page.route("**/wake-sessions/check-in", async (route) => {
    const body = route.request().postDataJSON();
    checkInRequests.push(body);
    await route.fulfill({
      status: 201,
      contentType: "application/json",
      body: JSON.stringify({
        wake_session_id: "session-123",
        status: "confirmed",
        streak: 4,
        checked_in_at: "2026-06-28T07:06:00.000Z",
        correlation_id: body.correlation_id
      })
    });
  });

  await page.route("**/wake-sessions/history**", async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        sessions: [
          {
            id: "session-123",
            status: "confirmed",
            target_time: "07:00",
            morning_intent: "Finish Kafka outbox implementation",
            checked_in_at: "2026-06-28T07:06:00.000Z",
            streak: 4
          }
        ]
      })
    });
  });

  await page.goto("/login");
  await page.getByLabel("Email").fill("marcus@example.com");
  await page.getByLabel("Password").fill("morning123");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page.getByText("Authenticated as demo-user")).toBeVisible();

  await page.goto("/wake-check-in");
  await page.getByRole("button", { name: "Validate desktop" }).click();
  await expect(page.getByText("Desktop validation passed. Morning Intent is now required.")).toBeVisible();

  await page.getByLabel("What is the most important thing you need to accomplish today?").fill(
    "Finish Kafka outbox implementation"
  );
  await page.getByRole("button", { name: "Start Session" }).click();
  await expect(page.getByText("Wake Session started and Morning Intent submitted.")).toBeVisible();
  await expect(page.getByText("Wake Session confirmed")).toBeVisible();
  await expect(page.getByText("Session session-123 | streak 4")).toBeVisible();

  await page.goto("/dashboard");
  await expect(page.getByText("Your morning in one glance.")).toBeVisible();
  await expect(page.getByText("75%")).toBeVisible();
  await expect(page.getByText("Current streak")).toBeVisible();
  await expect(page.getByText("Finish Kafka outbox implementation")).toBeVisible();

  expect(loginRequests).toHaveLength(1);
  expect(validationRequests).toHaveLength(1);
  expect(checkInRequests).toHaveLength(1);
  expect(validationRequests[0]).toMatchObject({
    user_agent: expect.any(String),
    viewport_width: expect.any(Number),
    viewport_height: expect.any(Number),
    screen_width: expect.any(Number),
    screen_height: expect.any(Number),
    has_touch: false,
    correlation_id: expect.any(String)
  });
  expect(checkInRequests[0]).toMatchObject({
    user_id: "demo-user",
    target_time: "07:00",
    morning_intent: "Finish Kafka outbox implementation",
    device_proof_id: "desktop-proof-123"
  });
});
