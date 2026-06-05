"use client";

import * as React from "react";
import { Trophy } from "lucide-react";

import { ProviderCard } from "@/components/ProviderCard";
import { TierToggle } from "@/components/TierToggle";
import {
  cheapest,
  flattenForTier,
  formatCost,
  providerLabel,
  type CostsResponse,
  type PlanEntry,
  type Tier,
} from "@/lib/plans";

interface ResultsViewProps {
  data: CostsResponse;
}

function groupByProvider(entries: PlanEntry[]): Map<string, PlanEntry[]> {
  const groups = new Map<string, PlanEntry[]>();
  for (const entry of entries) {
    const list = groups.get(entry.provider) ?? [];
    list.push(entry);
    groups.set(entry.provider, list);
  }
  return groups;
}

export function ResultsView({ data }: ResultsViewProps) {
  const [tier, setTier] = React.useState<Tier>("Standard");

  const entries = React.useMemo(() => flattenForTier(data, tier), [data, tier]);
  const best = cheapest(entries);
  const groups = groupByProvider(entries);

  return (
    <section className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h2 className="text-lg font-semibold">Estimated costs</h2>
        <TierToggle value={tier} onChange={setTier} />
      </div>

      {best && (
        <div className="flex items-center gap-4 rounded-xl bg-primary px-5 py-4 text-primary-foreground">
          <Trophy className="size-6 shrink-0" />
          <div className="flex-1">
            <p className="text-xs uppercase tracking-wide opacity-80">
              Cheapest plan
            </p>
            <p className="text-base font-semibold">
              {providerLabel(best.provider)} — {best.label}
            </p>
          </div>
          <p className="font-mono text-lg font-bold tabular-nums">
            {formatCost(best.cost)}
          </p>
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2">
        {[...groups.entries()].map(([provider, plans]) => (
          <ProviderCard
            key={provider}
            providerLabel={providerLabel(provider)}
            plans={plans}
            cheapestKey={best?.planKey ?? null}
          />
        ))}
      </div>

      <p className="text-xs text-muted-foreground">
        Costs are an estimate based on the usage data you uploaded.
      </p>
    </section>
  );
}
