"use client";

import { AppShell } from "@/components/AppShell";
import { checkIn, newCorrelationId, validateDevice } from "@/lib/api";
import { collectDeviceSignal, looksMobileLocally } from "@/lib/device";
import { CheckCircle2, Clock3, Compass, MonitorX, Send, Sparkles } from "lucide-react";
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
      <section className="mx-auto grid max-w-6xl gap-8 px-6 py-14 lg:grid-cols-[0.9fr_1.1fr] lg:py-20">
        <aside className="rounded-[32px] border border-black/5 bg-white/82 p-6 shadow-[0_20px_80px_rgba(15,23,42,0.08)] backdrop-blur-sm">
          <div className="inline-flex items-center gap-2 rounded-full border border-amber-200 bg-amber-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-amber-800">
            <Compass size={14} />
            Morning ritual
          </div>
          <h1 className="display-font mt-6 text-5xl leading-[0.95] text-slate-950">Good morning, Marcus.</h1>
          <p className="mt-4 max-w-sm text-base leading-7 text-slate-600">
            Take a breath, validate the desktop, and define the one thing that matters before the day gets loud.
          </p>
          <dl className="mt-8 grid gap-4 rounded-3xl border border-slate-100 bg-slate-50/80 p-5 text-sm">
            <div className="flex items-center justify-between border-b border-slate-200/70 pb-3">
              <dt className="text-slate-500">Planned wake-up</dt>
              <dd className="font-semibold text-slate-950">07:00</dd>
            </div>
            <div className="flex items-center justify-between border-b border-slate-200/70 pb-3">
              <dt className="text-slate-500">Current time</dt>
              <dd className="font-semibold text-slate-950">{now.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</dd>
            </div>
            <div className="flex items-center justify-between">
              <dt className="text-slate-500">Desktop proof</dt>
              <dd className="font-semibold text-slate-950">{deviceProofId ? "validated" : "pending"}</dd>
            </div>
          </dl>
          <button
            type="button"
            onClick={validateDesktop}
            className="focus-ring mt-6 inline-flex w-full items-center justify-center gap-2 rounded-full bg-slate-950 px-5 py-3 font-semibold text-white shadow-lg shadow-slate-950/10 transition hover:-translate-y-0.5"
          >
            {deviceProofId ? <CheckCircle2 size={18} /> : <MonitorX size={18} />}
            Validate desktop
          </button>
        </aside>

        <form onSubmit={submit} className="rounded-[32px] border border-black/5 bg-white/85 p-8 shadow-[0_20px_80px_rgba(15,23,42,0.08)] backdrop-blur-sm">
          <div className="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">
            <Sparkles size={14} />
            Morning Intent
          </div>
          <label className="display-font mt-5 block text-4xl leading-tight text-slate-950" htmlFor="intent">
            What is the most important thing you want to accomplish today?
          </label>
          <p className="mt-3 max-w-2xl text-base leading-7 text-slate-600">
            This is the heart of the session. Say the mission clearly, then start.
          </p>
          <textarea
            id="intent"
            value={intent}
            onChange={(event) => setIntent(event.target.value)}
            rows={8}
            className="focus-ring mt-6 w-full resize-none rounded-[28px] border border-slate-200 bg-slate-50/80 px-5 py-5 text-lg leading-8 text-slate-900 shadow-inner"
            placeholder="Finish Kafka outbox implementation"
          />
          <div className="mt-6 flex items-center justify-between gap-4 rounded-3xl border border-amber-200 bg-amber-50/80 px-5 py-4 text-sm text-amber-900">
            <div className="flex items-center gap-2">
              <Clock3 size={16} />
              A calm start matters more than a fast one.
            </div>
            <div className="hidden text-amber-800/80 sm:block">Desktop only MVP</div>
          </div>
          <button
            type="submit"
            className="focus-ring mt-6 inline-flex items-center gap-2 rounded-full px-5 py-3 font-semibold shadow-lg transition hover:-translate-y-0.5"
            style={{
              backgroundColor: "#0f172a",
              color: "#ffffff",
              border: "1px solid rgba(15, 23, 42, 0.14)",
              boxShadow: "0 14px 40px rgba(15, 23, 42, 0.16)"
            }}
          >
            <Send size={18} />
            Start Session
          </button>
          <p className="mt-4 text-sm text-slate-600">{status}</p>
          {result && (
            <div className="mt-6 rounded-3xl border border-emerald-200 bg-emerald-50 p-5">
              <div className="font-semibold text-emerald-900">Wake Session {result.status}</div>
              <div className="mt-1 text-sm text-emerald-800/80">Session {result.id} | streak {result.streak}</div>
            </div>
          )}
        </form>
      </section>
    </AppShell>
  );
}
