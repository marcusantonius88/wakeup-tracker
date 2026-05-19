"use client";

import { AppShell } from "@/components/AppShell";
import { checkIn, newCorrelationId, validateDevice } from "@/lib/api";
import { collectDeviceSignal, looksMobileLocally } from "@/lib/device";
import { CheckCircle2, MonitorX, Send } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

export default function WakeCheckInPage() {
  const [intent, setIntent] = useState("");
  const [interaction, setInteraction] = useState({ keyboard: false, pointer: false });
  const [deviceProofId, setDeviceProofId] = useState("");
  const [status, setStatus] = useState("Move the mouse or press a key, then validate this desktop.");
  const [result, setResult] = useState<{ id: string; streak: number; status: string } | null>(null);
  const correlationId = useMemo(() => newCorrelationId(), []);
  const now = new Date();

  useEffect(() => {
    function key() {
      setInteraction((value) => ({ ...value, keyboard: true }));
    }
    function pointer() {
      setInteraction((value) => ({ ...value, pointer: true }));
    }
    window.addEventListener("keydown", key);
    window.addEventListener("pointermove", pointer);
    window.addEventListener("pointerdown", pointer);
    return () => {
      window.removeEventListener("keydown", key);
      window.removeEventListener("pointermove", pointer);
      window.removeEventListener("pointerdown", pointer);
    };
  }, []);

  async function validateDesktop() {
    if (looksMobileLocally()) {
      setStatus("Mobile and tablet check-ins are blocked in the MVP.");
      return;
    }
    const payload = collectDeviceSignal(correlationId, interaction);
    const response = await validateDevice(payload);
    if (!response.validation.allowed) {
      setDeviceProofId("");
      setStatus(response.validation.reasons.join(" "));
      return;
    }
    setDeviceProofId(response.device_proof_id);
    setStatus("Desktop validation passed. Morning Intent is now required.");
  }

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!deviceProofId) {
      setStatus("Validate desktop access before starting the session.");
      return;
    }
    if (intent.trim().length < 8) {
      setStatus("Write a meaningful Morning Intent before confirming wake-up.");
      return;
    }

    try {
      const userID = localStorage.getItem("wakeup:user_id") ?? "demo-user";
      const response = await checkIn({
        user_id: userID,
        target_time: "07:00",
        checked_in_at: new Date().toISOString(),
        morning_intent: intent,
        device_proof_id: deviceProofId,
        correlation_id: correlationId
      });
      setResult({ id: response.wake_session_id, streak: response.streak, status: response.status });
      setStatus("Wake Session started and Morning Intent submitted.");
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Wake check-in failed");
    }
  }

  return (
    <AppShell>
      <section className="mx-auto grid max-w-6xl gap-6 px-6 py-10 lg:grid-cols-[0.9fr_1.1fr]">
        <aside className="rounded-md border border-line bg-white p-5">
          <h1 className="text-3xl font-semibold text-ink">Good morning Marcus</h1>
          <dl className="mt-6 grid gap-4 text-sm">
            <div className="flex items-center justify-between border-b border-line pb-3">
              <dt className="text-ink/60">Target wake-up time</dt>
              <dd className="font-semibold">07:00</dd>
            </div>
            <div className="flex items-center justify-between border-b border-line pb-3">
              <dt className="text-ink/60">Current time</dt>
              <dd className="font-semibold">{now.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</dd>
            </div>
            <div className="flex items-center justify-between">
              <dt className="text-ink/60">Desktop proof</dt>
              <dd className="font-semibold">{deviceProofId ? "validated" : "pending"}</dd>
            </div>
          </dl>
          <button
            type="button"
            onClick={validateDesktop}
            className="focus-ring mt-6 inline-flex w-full items-center justify-center gap-2 rounded-md border border-line bg-paper px-4 py-3 font-semibold text-ink"
          >
            {deviceProofId ? <CheckCircle2 size={18} /> : <MonitorX size={18} />}
            Validate desktop
          </button>
        </aside>

        <form onSubmit={submit} className="rounded-md border border-line bg-white p-5">
          <label className="text-lg font-semibold text-ink" htmlFor="intent">
            What is the most important thing you need to accomplish today?
          </label>
          <textarea
            id="intent"
            value={intent}
            onChange={(event) => setIntent(event.target.value)}
            rows={7}
            className="focus-ring mt-4 w-full resize-none rounded-md border border-line px-4 py-3 text-lg"
            placeholder="Finish Kafka outbox implementation"
          />
          <button className="focus-ring mt-5 inline-flex items-center gap-2 rounded-md bg-moss px-4 py-3 font-semibold text-paper">
            <Send size={18} />
            Start Session
          </button>
          <p className="mt-4 text-sm text-ink/65">{status}</p>
          {result && (
            <div className="mt-5 rounded-md border border-mint bg-mint/35 p-4">
              <div className="font-semibold text-ink">Wake Session {result.status}</div>
              <div className="mt-1 text-sm text-ink/70">Session {result.id} | streak {result.streak}</div>
            </div>
          )}
        </form>
      </section>
    </AppShell>
  );
}

