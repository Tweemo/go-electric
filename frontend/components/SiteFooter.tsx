import { Zap } from "lucide-react";

export function SiteFooter() {
  return (
    <footer className="bg-heading px-7 py-8">
      <div className="mx-auto flex w-full max-w-[1080px] flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <div className="flex size-7 items-center justify-center rounded-md bg-primary">
            <Zap size={14} className="fill-primary-foreground text-primary-foreground" />
          </div>
          <span className="text-base font-extrabold text-card">go-electric</span>
        </div>
        <p className="text-[13px] text-muted-foreground">
          Your file is processed in memory to calculate costs — it&apos;s never stored.
        </p>
      </div>
    </footer>
  );
}
