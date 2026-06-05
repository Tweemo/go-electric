"use client";

import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { TIERS, type Tier } from "@/lib/plans";

interface TierToggleProps {
  value: Tier;
  onChange: (tier: Tier) => void;
}

const TIER_DESCRIPTIONS: Record<Tier, string> = {
  Standard: "Standard user",
  Low: "Low user",
};

export function TierToggle({ value, onChange }: TierToggleProps) {
  return (
    <ToggleGroup
      value={[value]}
      onValueChange={(groupValue) => {
        // Single-select: keep the current value if the user deselects.
        const next = groupValue.find((v) => v !== value) ?? value;
        onChange(next as Tier);
      }}
      variant="outline"
      spacing={0}
    >
      {TIERS.map((tier) => (
        <ToggleGroupItem key={tier} value={tier} className="px-4">
          {TIER_DESCRIPTIONS[tier]}
        </ToggleGroupItem>
      ))}
    </ToggleGroup>
  );
}
