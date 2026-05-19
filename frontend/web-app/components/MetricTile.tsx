export function MetricTile({ label, value, detail }: { label: string; value: string; detail: string }) {
  return (
    <div className="rounded-md border border-line bg-white p-4">
      <div className="text-sm text-ink/65">{label}</div>
      <div className="mt-2 text-3xl font-semibold text-ink">{value}</div>
      <div className="mt-1 text-sm text-ink/60">{detail}</div>
    </div>
  );
}

