import { Zap } from "lucide-react";

interface SiteNavProps {
  onSample: () => void;
}

export function SiteNav({ onSample }: SiteNavProps) {
  return (
    <nav className="sticky top-0 z-40 border-b border-border/60 bg-background/90 backdrop-blur-md">
      <div className="mx-auto flex h-16 w-full max-w-[1080px] items-center justify-between px-7">
        <div className="flex items-center gap-2">
          <div className="flex size-8 items-center justify-center rounded-lg bg-primary">
            <Zap size={16} className="fill-primary-foreground text-primary-foreground" />
          </div>
          <span className="text-lg font-extrabold tracking-tight text-heading">
            go-electric
          </span>
        </div>
        <div className="flex items-center gap-6 text-sm">
          <a
            href="#how-section"
            className="font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            How it works
          </a>
          <button
            type="button"
            onClick={onSample}
            className="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary-dark"
          >
            Sample CSV
          </button>
        </div>
      </div>
    </nav>
  );
}
