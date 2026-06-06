import { Download, Upload, Zap } from "lucide-react";

const STEPS = [
  {
    icon: Download,
    step: "Step 1",
    title: "Export your usage",
    body: "Download your electricity usage history as a CSV from your provider's website or app — usually in the billing or usage section.",
  },
  {
    icon: Upload,
    step: "Step 2",
    title: "Upload the file",
    body: "Drop your CSV into the uploader. It's sent securely to our calculator, processed in memory, and never written to disk or stored.",
  },
  {
    icon: Zap,
    step: "Step 3",
    title: "See who's cheapest",
    body: "We calculate what each retailer's plans would charge based on your real usage and rank them from cheapest to priciest.",
  },
];

export function HowItWorks() {
  return (
    <section id="how-section" className="py-20">
      <div className="mx-auto w-full max-w-[1080px] px-7">
        <div className="mb-14 text-center">
          <h2 className="mb-3 text-[clamp(2rem,4vw,3rem)] font-extrabold tracking-tight text-heading">
            Three steps to savings.
          </h2>
          <p className="mx-auto max-w-md text-[17px] text-muted-foreground">
            No signup, no account, no upsells. Just your data and the answer.
          </p>
        </div>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {STEPS.map(({ icon: Icon, step, title, body }) => (
            <div
              key={step}
              className="rounded-2xl border border-border-subtle bg-card p-7 shadow-[0_2px_16px_rgba(20,38,30,0.04)]"
            >
              <div className="mb-5 flex size-12 items-center justify-center rounded-xl bg-accent">
                <Icon size={22} className="text-primary" />
              </div>
              <div className="mb-2 text-xs font-bold uppercase tracking-wider text-primary">
                {step}
              </div>
              <div className="mb-2.5 text-lg font-bold text-heading">{title}</div>
              <div className="text-sm leading-relaxed text-muted-foreground">
                {body}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
