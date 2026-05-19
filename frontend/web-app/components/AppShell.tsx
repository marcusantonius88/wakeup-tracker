import Link from "next/link";
import type { ReactNode } from "react";
import { AlarmClock, BarChart3, CheckCircle2, LogIn } from "lucide-react";

const nav = [
  { href: "/wake-check-in", label: "Check-in", icon: CheckCircle2 },
  { href: "/dashboard", label: "Dashboard", icon: BarChart3 },
  { href: "/login", label: "Login", icon: LogIn }
];

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <main className="min-h-screen">
      <header className="border-b border-line bg-paper/95">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <Link href="/" className="flex items-center gap-3 font-semibold text-ink">
            <span className="grid size-10 place-items-center rounded-md bg-moss text-paper">
              <AlarmClock size={20} />
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
                  className="focus-ring inline-flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium text-ink hover:bg-mint/50"
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
