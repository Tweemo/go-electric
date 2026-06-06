"use client";

import * as React from "react";

import { TierToggle } from "@/components/TierToggle";
import {
  cheapest,
  flattenForTier,
  formatCost,
  type CostsResponse,
  type Tier,
} from "@/lib/plans";

interface ResultsViewProps {
  data: CostsResponse;
}

export const ResultsView = React.forwardRef<HTMLElement, ResultsViewProps>(
  function ResultsView({ data }, ref) {
    const [tier, setTier] = React.useState<Tier>("Standard");

    const entries = React.useMemo(
      () => flattenForTier(data, tier),
      [data, tier],
    );
    const best = cheapest(entries);
    const mostExpensive =
      entries.length > 0 ? entries[entries.length - 1] : null;
    const maxCost = mostExpensive?.cost ?? 1;
    const runnersUp = entries.slice(1);

    return (
      <section
        ref={ref}
        className="scroll-mt-20 border-t border-border bg-card py-16"
      >
        <div className="mx-auto w-full max-w-[1080px] px-7">
          <div className="mb-8 flex flex-wrap items-center justify-between gap-3">
            <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
              Results — ranked cheapest first
            </h2>
            <TierToggle value={tier} onChange={setTier} />
          </div>

          {best && (
            <>
              {/* Winner card */}
              <div className="relative mb-4 overflow-hidden rounded-2xl bg-gradient-to-br from-primary to-primary-dark p-8 text-primary-foreground">
                <div className="pointer-events-none absolute -top-10 -right-10 size-52 rounded-full bg-white/5" />
                <div className="pointer-events-none absolute top-5 right-5 size-24 rounded-full bg-white/[0.06]" />
                <div className="relative">
                  <div className="mb-4 inline-flex items-center gap-1.5 rounded-full bg-white/20 px-3.5 py-1 text-xs font-bold tracking-wide">
                    🏆 Cheapest for your usage
                  </div>
                  <div className="flex flex-wrap items-start justify-between gap-6">
                    <div>
                      <div className="text-[clamp(1.6rem,3.5vw,2.5rem)] font-extrabold leading-tight tracking-tight">
                        {best.providerLabel}
                      </div>
                      <div className="text-[15px] opacity-80">{best.label}</div>
                    </div>
                    <div className="text-right">
                      <div className="mb-1 text-[13px] opacity-75">
                        Estimated cost for this period
                      </div>
                      <div className="text-[clamp(2.5rem,5vw,4rem)] font-extrabold leading-none tracking-tight">
                        {formatCost(best.cost)}
                      </div>
                    </div>
                  </div>
                  <div className="mt-7 grid gap-4 border-t border-white/20 pt-6 sm:grid-cols-3">
                    {[
                      { label: "Provider", value: best.providerLabel },
                      { label: "Plan", value: best.label },
                      {
                        label: "You could save",
                        value: mostExpensive
                          ? formatCost(mostExpensive.cost - best.cost)
                          : "—",
                      },
                    ].map((s) => (
                      <div key={s.label}>
                        <div className="mb-1 text-xs opacity-70">{s.label}</div>
                        <div className="text-base font-bold">{s.value}</div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Ranked runners-up */}
              <div className="flex flex-col gap-2">
                {runnersUp.map((c, idx) => {
                  const diff = c.cost - best.cost;
                  const barPct = (c.cost / maxCost) * 100;
                  const winnerPct = (best.cost / maxCost) * 100;
                  return (
                    <div
                      key={c.planKey}
                      className="flex flex-wrap items-center gap-5 rounded-xl border border-border-subtle bg-muted px-6 py-4.5"
                    >
                      <span className="w-5 text-sm font-semibold text-subtle">
                        #{idx + 2}
                      </span>
                      <div className="min-w-[200px] flex-1">
                        <div className="mb-2.5 flex flex-wrap items-center justify-between gap-4">
                          <div>
                            <span className="text-base font-bold text-heading">
                              {c.providerLabel}
                            </span>
                            <span className="ml-2.5 text-xs font-medium text-subtle">
                              {c.label}
                            </span>
                          </div>
                          <span className="whitespace-nowrap rounded-full bg-destructive-soft px-2.5 py-0.5 text-[13px] font-semibold text-destructive">
                            +{formatCost(diff)} more
                          </span>
                        </div>
                        <div className="h-1.5 overflow-hidden rounded-full bg-border-subtle">
                          <div
                            className="relative h-full rounded-full bg-border"
                            style={{ width: `${barPct}%` }}
                          >
                            <div
                              className="absolute top-0 left-0 h-full rounded-full bg-primary"
                              style={{ width: `${(winnerPct / barPct) * 100}%` }}
                            />
                          </div>
                        </div>
                      </div>
                      <div className="min-w-20 text-right">
                        <div className="font-mono text-lg font-extrabold tabular-nums text-heading">
                          {formatCost(c.cost)}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </>
          )}

          <p className="mt-5 text-xs leading-relaxed text-subtle">
            * Costs are estimated from the usage data you uploaded and cover the
            whole period in your file. Prices are indicative only — contact
            retailers for current accurate pricing.
          </p>
        </div>
      </section>
    );
  },
);
