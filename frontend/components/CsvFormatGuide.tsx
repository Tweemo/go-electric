import { Download } from "lucide-react";

interface CsvFormatGuideProps {
  onSample: () => void;
}

const COLUMNS = [
  {
    label: "Start & end date-time",
    note: "When each reading interval begins and ends, e.g. 09/03/2025 00:00:00.",
  },
  {
    label: "Usage value",
    note: "Energy used in that interval, in kWh — e.g. 0.775.",
  },
];

export function CsvFormatGuide({ onSample }: CsvFormatGuideProps) {
  return (
    <section className="bg-card py-16">
      <div className="mx-auto w-full max-w-[1080px] px-7">
        <div className="grid items-start gap-12 md:grid-cols-2">
          <div>
            <h3 className="mb-3.5 text-[28px] font-extrabold tracking-tight text-heading">
              CSV format guide
            </h3>
            <p className="mb-7 text-[15px] leading-relaxed text-muted-foreground">
              We read the standard half-hourly meter export (the ICP consumption
              file) most NZ retailers provide. Each <code className="rounded bg-muted px-1 font-mono text-[13px]">DET</code>{" "}
              row is one reading interval; we use the date-time and usage columns.
            </p>
            <div className="mb-7 flex flex-col gap-3">
              {COLUMNS.map((c) => (
                <div
                  key={c.label}
                  className="rounded-xl border border-border bg-muted px-4.5 py-4"
                >
                  <div className="mb-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                    {c.label}
                  </div>
                  <div className="text-sm text-foreground">{c.note}</div>
                </div>
              ))}
            </div>
            <button
              type="button"
              onClick={onSample}
              className="flex items-center gap-2 rounded-[10px] border-[1.5px] border-border bg-muted px-5 py-3 text-sm font-semibold text-primary transition-colors hover:bg-accent"
            >
              <Download size={15} /> Download sample CSV
            </button>
          </div>
          <div>
            <div className="mb-3 text-xs font-bold uppercase tracking-wide text-subtle">
              Example file
            </div>
            <div className="overflow-x-auto rounded-2xl bg-code px-6 py-5 font-mono text-[13px] leading-7">
              <div className="mb-1 text-muted-foreground"># sample-usage.csv</div>
              <div className="text-subtle">HDR,ICPCONS,1.1,…</div>
              <div className="whitespace-nowrap text-code-foreground">
                DET,…,09/03/2025 00:00:00,09/03/2025 00:30:00,RD,0.775
              </div>
              <div className="whitespace-nowrap text-code-foreground">
                DET,…,09/03/2025 00:30:00,09/03/2025 01:00:00,RD,0.241
              </div>
              <div className="whitespace-nowrap text-code-foreground">
                DET,…,09/03/2025 01:00:00,09/03/2025 01:30:00,RD,0.198
              </div>
              <div className="text-subtle">… more rows</div>
            </div>
            <p className="mt-3 text-[13px] text-subtle">
              Half-hourly readings across the whole period you upload are totalled —
              a six-month file gives a six-month cost.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
