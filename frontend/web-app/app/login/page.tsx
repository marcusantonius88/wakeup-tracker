"use client";

import { AppShell } from "@/components/AppShell";
import { login } from "@/lib/api";
import { LogIn } from "lucide-react";
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
      <section className="mx-auto max-w-xl px-6 py-10">
        <h1 className="text-3xl font-semibold text-ink">Login</h1>
        <form onSubmit={submit} className="mt-6 rounded-md border border-line bg-white p-5">
          <label className="block text-sm font-medium text-ink" htmlFor="email">Email</label>
          <input
            id="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            className="focus-ring mt-2 w-full rounded-md border border-line px-3 py-3"
          />
          <label className="mt-4 block text-sm font-medium text-ink" htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            className="focus-ring mt-2 w-full rounded-md border border-line px-3 py-3"
          />
          <button className="focus-ring mt-5 inline-flex items-center gap-2 rounded-md bg-moss px-4 py-3 font-semibold text-paper">
            <LogIn size={18} />
            Sign in
          </button>
          <p className="mt-4 text-sm text-ink/65">{status}</p>
        </form>
      </section>
    </AppShell>
  );
}

