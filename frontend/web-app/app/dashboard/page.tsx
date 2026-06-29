"use client";

import { AppShell } from "@/components/AppShell";
import { MetricTile } from "@/components/MetricTile";
import { loadHistory } from "@/lib/api";
import { Activity, CalendarClock, Target, ArrowUpRight, Sparkles } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

type WakeSession = {
  id: string;
  status: string;
  target_time: string;
  morning_intent: string;
  checked_in_at: string;
  streak: number;
};

export default function DashboardPage() {
  const [sessions, setSessions] = useState<WakeSession[]>([]);

  useEffect(() => {
    const userID = localStorage.getItem("wakeup:user_id") ?? "demo-user";
    loadHistory(userID).then((data) => setSessions(data.sessions ?? []));
  }, []);

  const confirmed = sessions.filter((session) => session.status === "confirmed").length;
  const failed = sessions.filter((session) => session.status === "failed").length;
  const consistency = useMemo(() => {
    const total = confirmed + failed;
    return total === 0 ? 0 : Math.round((confirmed / total) * 100);
  }, [confirmed, failed]);
  const streak = sessions.at(-1)?.streak ?? 0;

  return (
    <AppShell>
      <section className="mx-auto max-w-6xl px-6 py-14 lg:py-20">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <div className="inline-flex items-center gap-2 rounded-full border border-black/5 bg-white/80 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">
              <Sparkles size={14} />
              Post-session insight
            </div>
            <h1 className="display-font mt-4 text-5xl leading-none text-slate-950">Your morning in one glance.</h1>
            <p className="mt-4 max-w-2xl text-base leading-7 text-slate-600">
              A clean read model of consistency, streak, average wake-up, and recent Morning Intents without the noise of a heavy dashboard.
            </p>
          </div>
          <div className="rounded-full border border-amber-200 bg-amber-50 px-4 py-2 text-sm font-semibold text-amber-800">
            Desktop-first MVP
          </div>
        </div>

        <div className="mt-8 grid gap-4 md:grid-cols-4">
          <MetricTile label="Weekly consistency" value={`${consistency}%`} detail="WakeUpConfirmed over total sessions" />
          <MetricTile label="Current streak" value={`${streak}`} detail="Updated by StreakIncreased events" />
          <MetricTile label="Success count" value={`${confirmed}`} detail="Confirmed wake-ups" />
          <MetricTile label="Morning Intents" value={`${sessions.length}`} detail="Daily cognitive activation" />
        </div>

        <div className="mt-8 grid gap-6 lg:grid-cols-[1fr_0.8fr]">
          <section className="overflow-hidden rounded-[32px] border border-black/5 bg-white/85 shadow-[0_20px_80px_rgba(15,23,42,0.08)]">
            <div className="flex items-center gap-2 border-b border-slate-100 px-6 py-4 font-semibold text-slate-900">
              <Activity size={18} />
              Wake-up timeline
            </div>
            <div className="divide-y divide-slate-100">
              {sessions.length === 0 && <p className="px-6 py-8 text-sm text-slate-500">No sessions yet.</p>}
              {sessions.map((session) => (
                <div key={session.id} className="grid gap-1 px-6 py-4 transition hover:bg-slate-50/80">
                  <div className="flex items-center justify-between gap-3">
                    <span className="font-semibold text-slate-950">{session.status}</span>
                    <span className="text-sm text-slate-500">{new Date(session.checked_in_at).toLocaleString()}</span>
                  </div>
                  <p className="text-sm leading-6 text-slate-600">{session.morning_intent}</p>
                </div>
              ))}
            </div>
          </section>

          <section className="rounded-[32px] border border-black/5 bg-white/85 shadow-[0_20px_80px_rgba(15,23,42,0.08)]">
            <div className="flex items-center gap-2 border-b border-slate-100 px-6 py-4 font-semibold text-slate-900">
              <Target size={18} />
              Behavioral signals
            </div>
            <div className="grid gap-4 p-6">
              <div className="flex items-center gap-3 rounded-3xl border border-slate-100 bg-slate-50/80 p-4">
                <CalendarClock className="text-amber-500" size={22} />
                <div>
                  <div className="font-semibold text-slate-950">Average wake-up time</div>
                  <div className="text-sm text-slate-500">Projected once multiple sessions exist</div>
                </div>
              </div>
              <div className="rounded-3xl border border-slate-100 p-5">
                <div className="flex items-center justify-between gap-3">
                  <div className="font-semibold text-slate-950">Morning Intent pattern</div>
                  <ArrowUpRight size={18} className="text-slate-400" />
                </div>
                <p className="mt-3 text-sm leading-7 text-slate-600">
                  Future analytics can classify intent themes, recurrence, and goal completion against productivity sessions.
                </p>
              </div>
            </div>
          </section>
        </div>
      </section>
    </AppShell>
  );
}
