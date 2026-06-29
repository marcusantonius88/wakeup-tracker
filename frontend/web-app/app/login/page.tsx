"use client";

import { AppShell } from "@/components/AppShell";
import { login } from "@/lib/api";
import { ArrowRight, LogIn, Mail, Shield } from "lucide-react";
import { useState } from "react";

export default function LoginPage() {
  const [email, setEmail] = useState("marcus@example.com");
  const [password, setPassword] = useState("morning123");
  const [status, setStatus] = useState("Ready");

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setStatus("Authenticating...");
    try {
      const response = await login(email, password);
      localStorage.setItem("wakeup:user_id", response.user_id);
      localStorage.setItem("wakeup:access_token", response.access_token);
      setStatus(`Authenticated as ${response.user_id}`);
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Login failed");
    }
  }

  return (
    <AppShell>
      <section className="mx-auto max-w-xl px-6 py-16">
        <div className="rounded-[32px] border border-black/5 bg-white/80 p-8 shadow-[0_20px_80px_rgba(15,23,42,0.08)] backdrop-blur-sm">
          <div className="inline-flex items-center gap-2 rounded-full border border-amber-200 bg-amber-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-amber-800">
            <Shield size={14} />
            Secure entry
          </div>
          <h1 className="display-font mt-5 text-5xl leading-none text-slate-950">Welcome back</h1>
          <p className="mt-4 max-w-md text-base leading-7 text-slate-600">
            Sign in to continue your morning ritual and resume the workspace that helps you begin with intention.
          </p>
          <form onSubmit={submit} className="mt-8 grid gap-5">
            <div>
              <label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="email">Email</label>
              <div className="relative">
                <Mail className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                <input
                  id="email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  className="focus-ring w-full rounded-2xl border border-slate-200 bg-white px-11 py-3 text-slate-900 shadow-sm"
                />
              </div>
            </div>
            <div>
              <label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="password">Password</label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                className="focus-ring w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-slate-900 shadow-sm"
              />
            </div>
            <button className="focus-ring inline-flex w-fit items-center gap-2 rounded-full bg-slate-950 px-5 py-3 font-semibold text-white shadow-lg shadow-slate-950/10 transition hover:-translate-y-0.5">
              <LogIn size={18} />
              Sign in
              <ArrowRight size={16} />
            </button>
            <p className="text-sm text-slate-600">{status}</p>
          </form>
        </div>
      </section>
    </AppShell>
  );
}
