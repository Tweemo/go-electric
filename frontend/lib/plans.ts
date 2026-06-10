// Types and helpers for the cost-comparison response returned by the Go API.
// The API returns: { provider: { planKey: cost } }, e.g.
// { "contact": { "GoodChargeStandard": 1234.5, ... }, "nova": { ... } }

export type CostsResponse = Record<string, Record<string, number>>;

export type Tier = "Standard" | "Low";

export const TIERS: Tier[] = ["Standard", "Low"];

export const PROVIDER_LABELS: Record<string, string> = {
  contact: "Contact Energy",
  nova: "Nova Energy",
  powershop: "Powershop",
};

// Plan keys carry a tier suffix ("Standard" | "Low"). These labels cover the
// base plan name; tierOf() extracts the suffix separately.
const PLAN_LABELS: Record<string, string> = {
  GoodCharge: "Good Charge",
  GoodNights: "Good Nights",
  GoodWeekends: "Good Weekends",
  SimpleRates: "Simple Rates",
  GeneralRates: "General Rates",
  Basic: "Basic",
};

export function tierOf(planKey: string): Tier | null {
  if (planKey.endsWith("Standard")) return "Standard";
  if (planKey.endsWith("Low")) return "Low";
  return null;
}

export function planLabel(planKey: string): string {
  const tier = tierOf(planKey);
  const base = tier ? planKey.slice(0, planKey.length - tier.length) : planKey;
  return PLAN_LABELS[base] ?? base;
}

export function providerLabel(provider: string): string {
  return PROVIDER_LABELS[provider] ?? provider;
}

export interface PlanEntry {
  provider: string;
  providerLabel: string;
  planKey: string;
  label: string;
  cost: number;
}

// Flatten the nested response into a list of plans for one tier, cheapest first.
export function flattenForTier(resp: CostsResponse, tier: Tier): PlanEntry[] {
  const entries: PlanEntry[] = [];
  for (const [provider, plans] of Object.entries(resp)) {
    for (const [planKey, cost] of Object.entries(plans)) {
      if (tierOf(planKey) !== tier) continue;
      entries.push({
        provider,
        providerLabel: providerLabel(provider),
        planKey,
        label: planLabel(planKey),
        cost,
      });
    }
  }
  return entries.sort((a, b) => a.cost - b.cost);
}

export function cheapest(entries: PlanEntry[]): PlanEntry | null {
  return entries.length > 0 ? entries[0] : null;
}

const nzd = new Intl.NumberFormat("en-NZ", {
  style: "currency",
  currency: "NZD",
});

export function formatCost(cost: number): string {
  return nzd.format(cost);
}
