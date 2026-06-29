import Link from "next/link";
import type { ReactNode } from "react";
import { ArrowRight, BarChart3, CheckCircle2, LogIn, Sunrise } from "lucide-react";

const nav = [
  { href: "/wake-check-in", label: "Check-in", icon: CheckCircle2 },
  { href: "/dashboard", label: "Dashboard", icon: BarChart3 },
  { href: "/login", label: "Login", icon: LogIn }
];

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(255,223,180,0.35),_transparent_30%),linear-gradient(180deg,_#fcfbf7_0%,_#f5f1e8_100%)]">
      <header className="border-b border-black/5 bg-white/70 backdrop-blur-xl">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <Link href="/" className="flex items-center gap-3 text-sm font-semibold tracking-wide text-slate-900">
            <span className="grid size-10 place-items-center rounded-2xl bg-slate-900 text-white shadow-sm">
              <Sunrise size={18} />
            </span>
            WakeUpTracker
          </Link>
          <nav className="flex items-center gap-2">
            {nav.map((item) => {
              const Icon = item.icon;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className="focus-ring inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white/80 px-4 py-2 text-sm font-medium text-slate-700 shadow-sm transition hover:-translate-y-0.5 hover:border-slate-300 hover:text-slate-950"
                >
                  <Icon size={16} />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>
      </header>
      {children}
    </main>
  );
}
