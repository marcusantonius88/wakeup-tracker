import { expect, test } from "@playwright/test";

test("mobile check-ins are blocked in the MVP", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });

  await page.goto("/wake-check-in");
  await page.getByRole("button", { name: "Validate desktop" }).click();

  await expect(page.getByText("Mobile and tablet check-ins are blocked in the MVP.")).toBeVisible();
});

