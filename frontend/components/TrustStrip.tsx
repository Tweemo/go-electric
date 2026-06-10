import { Shield, Cpu, Leaf } from "lucide-react";

const ITEMS = [
  {
    icon: Shield,
    title: "Private by design",
    body: "Your ICP, meter number and name are stripped in your browser — only anonymous date and usage figures are sent, and nothing is ever stored.",
  },
  {
    icon: Cpu,
    title: "Based on real data",
    body: "Results use your actual usage readings, not estimates or averages.",
  },
  {
    icon: Leaf,
    title: "Completely free",
    body: "No account, no credit card, no upsells. Ever.",
  },
];

export function TrustStrip() {
  return (
    <section className="border-y border-border bg-accent px-7 py-10">
      <div className="mx-auto grid w-full max-w-[1080px] gap-8 sm:grid-cols-2 lg:grid-cols-3">
        {ITEMS.map(({ icon: Icon, title, body }) => (
          <div key={title} className="flex items-start gap-3.5">
            <div className="flex size-10 shrink-0 items-center justify-center rounded-[10px] bg-card">
              <Icon size={20} className="text-primary" />
            </div>
            <div>
              <div className="mb-1 text-[15px] font-bold text-heading">{title}</div>
              <div className="text-[13px] leading-snug text-muted-foreground">
                {body}
              </div>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
