import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { formatCost, type PlanEntry } from "@/lib/plans";

interface ProviderCardProps {
  providerLabel: string;
  plans: PlanEntry[];
  cheapestKey: string | null;
}

export function ProviderCard({
  providerLabel,
  plans,
  cheapestKey,
}: ProviderCardProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{providerLabel}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-1">
        {plans.map((plan) => {
          const isCheapest = plan.planKey === cheapestKey;
          return (
            <div
              key={plan.planKey}
              className={cn(
                "flex items-center justify-between rounded-lg px-3 py-2",
                isCheapest && "bg-primary/5 ring-1 ring-primary/20",
              )}
            >
              <span className="flex items-center gap-2 text-sm">
                {plan.label}
                {isCheapest && (
                  <Badge variant="secondary" className="text-xs">
                    Cheapest
                  </Badge>
                )}
              </span>
              <span className="font-mono text-sm font-medium tabular-nums">
                {formatCost(plan.cost)}
              </span>
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
