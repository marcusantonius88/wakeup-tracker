export function MetricTile({ label, value, detail }: { label: string; value: string; detail: string }) {
  return (
    <div className="rounded-3xl border border-black/5 bg-white/85 p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
      <div className="text-xs uppercase tracking-[0.22em] text-slate-500">{label}</div>
      <div className="mt-3 display-font text-4xl leading-none text-slate-950">{value}</div>
      <div className="mt-3 text-sm leading-6 text-slate-600">{detail}</div>
    </div>
  );
}
