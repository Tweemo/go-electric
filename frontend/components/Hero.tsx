import { Upload, Download, Leaf } from "lucide-react";

interface HeroProps {
  onUpload: () => void;
  onSample: () => void;
}

const EXAMPLE_RUNNERS_UP = [
  { name: "Nova Energy — General Rates", cost: 1923, diff: 76 },
  { name: "Contact Energy — Good Nights", cost: 1961, diff: 114 },
  { name: "Contact Energy — Good Weekends", cost: 2044, diff: 197 },
];

export function Hero({ onUpload, onSample }: HeroProps) {
  return (
    <section className="mx-auto w-full max-w-[1080px] px-7 pt-20 pb-12">
      <div className="grid items-center gap-16 md:grid-cols-2">
        {/* Copy */}
        <div>
          <div className="mb-6 inline-flex items-center gap-1.5 rounded-full bg-accent px-3.5 py-1.5 text-xs font-semibold text-accent-foreground">
            <Leaf size={12} /> Free · Private · Instant
          </div>
          <h1 className="mb-5 text-[clamp(2.4rem,5vw,4rem)] font-extrabold leading-[1.1] tracking-tight text-heading">
            Stop guessing.
            <br />
            <span className="text-primary">Start saving</span> on
            <br />
            electricity.
          </h1>
          <p className="mb-8 max-w-md text-[17px] leading-relaxed text-muted-foreground">
            Upload your power usage CSV and we&apos;ll instantly show you which
            retailer would give you the best deal — based on your real usage, not
            estimates or averages.
          </p>
          <div className="flex flex-wrap gap-3">
            <button
              type="button"
              onClick={onUpload}
              className="flex items-center gap-2 rounded-xl bg-primary px-6 py-3.5 text-base font-bold text-primary-foreground transition-colors hover:bg-primary-dark"
            >
              <Upload size={18} /> Upload my CSV
            </button>
            <button
              type="button"
              onClick={onSample}
              className="flex items-center gap-2 rounded-xl border-[1.5px] border-border bg-card px-6 py-3.5 text-[15px] font-semibold text-foreground transition-colors hover:bg-secondary"
            >
              <Download size={16} className="text-primary" /> Try a sample
            </button>
          </div>
        </div>

        {/* Preview card */}
        <div className="rounded-2xl border border-border-subtle bg-card p-6 shadow-[0_4px_40px_rgba(20,38,30,0.08)]">
          <div className="mb-4 flex items-center justify-between">
            <span className="text-[13px] font-semibold text-muted-foreground">
              Example result
            </span>
            <span className="rounded-full bg-muted px-2.5 py-1 text-xs text-subtle">
              182 days · 4,058 kWh
            </span>
          </div>
          <div className="mb-2.5 rounded-2xl border-2 border-primary bg-accent p-4">
            <div className="mb-1.5 text-[11px] font-bold tracking-wide text-accent-foreground">
              🏆 CHEAPEST FOR YOUR USAGE
            </div>
            <div className="mb-0.5 text-[22px] font-extrabold text-heading">
              Contact Energy — Good Charge
            </div>
            <div className="flex items-baseline gap-1.5">
              <span className="text-[32px] font-extrabold tracking-tight text-primary">
                $1,847
              </span>
              <span className="text-sm text-muted-foreground">/period</span>
            </div>
          </div>
          {EXAMPLE_RUNNERS_UP.map((r) => (
            <div
              key={r.name}
              className="flex items-center justify-between border-b border-border-subtle py-2.5 last:border-b-0"
            >
              <span className="text-sm font-medium text-foreground">{r.name}</span>
              <div className="flex items-center gap-2.5">
                <span className="text-sm font-semibold">
                  ${r.cost.toLocaleString()}
                </span>
                <span className="rounded-full bg-destructive-soft px-2 py-0.5 text-xs font-medium text-destructive">
                  +${r.diff}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
