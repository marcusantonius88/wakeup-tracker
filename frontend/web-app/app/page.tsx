import { AppShell } from "@/components/AppShell";
import { MetricTile } from "@/components/MetricTile";
import Link from "next/link";
import { ArrowRight, MonitorCheck } from "lucide-react";

export default function Home() {
  return (
    <AppShell>
      <section className="mx-auto grid max-w-6xl gap-8 px-6 py-10 lg:grid-cols-[1.1fr_0.9fr]">
        <div className="flex flex-col justify-center">
          <div className="mb-4 inline-flex w-fit items-center gap-2 rounded-md border border-line bg-white px-3 py-2 text-sm font-medium text-moss">
            <MonitorCheck size={16} />
            Desktop wake accountability
          </div>
          <h1 className="max-w-3xl text-5xl font-semibold tracking-normal text-ink">
            WakeUpTracker
          </h1>
          <p className="mt-5 max-w-2xl text-lg leading-8 text-ink/70">
            Confirm your wake-up only after opening the notebook, passing desktop validation,
            and writing the most important objective of the day.
          </p>
          <div className="mt-8 flex gap-3">
            <Link
              href="/wake-check-in"
              className="focus-ring inline-flex items-center gap-2 rounded-md bg-moss px-4 py-3 font-semibold text-paper"
            >
              Start wake check-in
              <ArrowRight size={18} />
            </Link>
            <Link
              href="/dashboard"
              className="focus-ring inline-flex items-center rounded-md border border-line bg-white px-4 py-3 font-semibold text-ink"
            >
              View dashboard
            </Link>
          </div>
        </div>
        <div className="grid content-start gap-4">
          <MetricTile label="Target wake-up" value="07:00" detail="10 minute tolerance window" />
          <MetricTile label="Current streak" value="3 days" detail="Projected from WakeUpConfirmed events" />
          <MetricTile label="Morning Intent" value="Required" detail="No passive check-ins in the MVP" />
        </div>
      </section>
    </AppShell>
  );
}

