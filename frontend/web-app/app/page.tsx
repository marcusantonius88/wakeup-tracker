import { AppShell } from "@/components/AppShell";
import { MetricTile } from "@/components/MetricTile";
import Link from "next/link";
import { ArrowRight, ChevronRight, MonitorCheck, Sparkles } from "lucide-react";

export default function Home() {
  return (
    <AppShell>
      <section className="mx-auto grid max-w-6xl gap-8 px-6 py-14 lg:grid-cols-[1.15fr_0.85fr] lg:py-20">
        <div className="flex flex-col justify-center">
          <div className="mb-6 inline-flex w-fit items-center gap-2 rounded-full border border-black/5 bg-white/80 px-4 py-2 text-sm font-medium text-slate-600 shadow-sm">
            <Sparkles size={16} className="text-amber-500" />
            Morning ritual for remote work
          </div>
          <h1 className="display-font max-w-3xl text-6xl leading-none tracking-tight text-slate-950 lg:text-7xl">
            Start the day with purpose.
          </h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-600">
            WakeUpTracker turns wake-up into a quiet ritual: confirm your desktop, define
            what matters, and begin the day with intention instead of autopilot.
          </p>
          <div className="mt-10 flex flex-wrap gap-3">
            <Link
              href="/wake-check-in"
              className="focus-ring inline-flex items-center gap-2 rounded-full bg-slate-950 px-5 py-3 font-semibold text-white shadow-lg shadow-slate-950/10 transition hover:-translate-y-0.5"
            >
              Start wake check-in
              <ArrowRight size={18} />
            </Link>
            <Link
              href="/dashboard"
              className="focus-ring inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white/80 px-5 py-3 font-semibold text-slate-700 shadow-sm transition hover:-translate-y-0.5 hover:border-slate-300"
            >
              View dashboard
              <ChevronRight size={18} />
            </Link>
          </div>
        </div>
        <div className="grid content-start gap-4 rounded-[28px] border border-black/5 bg-white/70 p-5 shadow-[0_20px_80px_rgba(15,23,42,0.08)] backdrop-blur-sm">
          <div className="rounded-3xl border border-black/5 bg-gradient-to-br from-slate-950 to-slate-800 p-6 text-white shadow-lg">
            <div className="text-sm uppercase tracking-[0.24em] text-white/60">Wake Check-in</div>
            <div className="mt-3 display-font text-3xl leading-tight">Your morning, centered.</div>
            <p className="mt-4 max-w-sm text-sm leading-6 text-white/70">
              Desktop validation first. Morning Intent second. A calm start, then a productive rhythm.
            </p>
          </div>
          <MetricTile label="Target wake-up" value="07:00" detail="10 minute tolerance window" />
          <MetricTile label="Current streak" value="3 days" detail="Projected from WakeUpConfirmed events" />
          <MetricTile label="Morning Intent" value="Required" detail="No passive check-ins in the MVP" />
        </div>
      </section>
    </AppShell>
  );
}
