"use client";

import { AppShell } from "@/components/AppShell";
import { MetricTile } from "@/components/MetricTile";
import { loadHistory } from "@/lib/api";
import { Activity, CalendarClock, Target } from "lucide-react";
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
      <section className="mx-auto max-w-6xl px-6 py-10">
        <div className="flex items-start justify-between gap-6">
          <div>
            <h1 className="text-3xl font-semibold text-ink">Consistency Dashboard</h1>
            <p className="mt-2 text-ink/65">CQRS read model for wake sessions, streaks, timeline, and Morning Intents.</p>
          </div>
          <div className="rounded-md border border-line bg-white px-4 py-3 text-sm font-semibold text-moss">
            MVP desktop-only
          </div>
        </div>

        <div className="mt-6 grid gap-4 md:grid-cols-4">
          <MetricTile label="Weekly consistency" value={`${consistency}%`} detail="WakeUpConfirmed over total sessions" />
          <MetricTile label="Current streak" value={`${streak}`} detail="Updated by StreakIncreased events" />
          <MetricTile label="Success count" value={`${confirmed}`} detail="Confirmed wake-ups" />
          <MetricTile label="Morning Intents" value={`${sessions.length}`} detail="Daily cognitive activation" />
        </div>

        <div className="mt-8 grid gap-6 lg:grid-cols-[1fr_0.8fr]">
          <section className="rounded-md border border-line bg-white">
            <div className="flex items-center gap-2 border-b border-line px-5 py-4 font-semibold">
              <Activity size={18} />
              Wake-up timeline
            </div>
            <div className="divide-y divide-line">
              {sessions.length === 0 && <p className="px-5 py-6 text-sm text-ink/60">No sessions yet.</p>}
              {sessions.map((session) => (
                <div key={session.id} className="grid gap-1 px-5 py-4">
                  <div className="flex items-center justify-between gap-3">
                    <span className="font-semibold text-ink">{session.status}</span>
                    <span className="text-sm text-ink/55">{new Date(session.checked_in_at).toLocaleString()}</span>
                  </div>
                  <p className="text-sm text-ink/65">{session.morning_intent}</p>
                </div>
              ))}
            </div>
          </section>

          <section className="rounded-md border border-line bg-white">
            <div className="flex items-center gap-2 border-b border-line px-5 py-4 font-semibold">
              <Target size={18} />
              Behavioral signals
            </div>
            <div className="grid gap-4 p-5">
              <div className="flex items-center gap-3 rounded-md border border-line p-4">
                <CalendarClock className="text-amber" size={22} />
                <div>
                  <div className="font-semibold">Average wake-up time</div>
                  <div className="text-sm text-ink/60">Projected once multiple sessions exist</div>
                </div>
              </div>
              <div className="rounded-md border border-line p-4">
                <div className="font-semibold">Morning Intent pattern</div>
                <p className="mt-2 text-sm leading-6 text-ink/65">
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

