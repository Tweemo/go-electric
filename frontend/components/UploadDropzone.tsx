"use client";

import * as React from "react";
import { Upload, FileText, Loader2 } from "lucide-react";

import { cn } from "@/lib/utils";

interface UploadDropzoneProps {
  onFile: (file: File) => void;
  loading: boolean;
  fileName: string | null;
}

export function UploadDropzone({
  onFile,
  loading,
  fileName,
}: UploadDropzoneProps) {
  const inputRef = React.useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = React.useState(false);

  function handleFiles(files: FileList | null) {
    const file = files?.[0];
    if (file) onFile(file);
  }

  return (
    <div
      role="button"
      tabIndex={0}
      aria-disabled={loading}
      onClick={() => !loading && inputRef.current?.click()}
      onKeyDown={(e) => {
        if ((e.key === "Enter" || e.key === " ") && !loading) {
          e.preventDefault();
          inputRef.current?.click();
        }
      }}
      onDragOver={(e) => {
        e.preventDefault();
        if (!loading) setDragging(true);
      }}
      onDragLeave={() => setDragging(false)}
      onDrop={(e) => {
        e.preventDefault();
        setDragging(false);
        if (!loading) handleFiles(e.dataTransfer.files);
      }}
      className={cn(
        "flex w-full cursor-pointer flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed px-6 py-12 text-center transition-colors",
        dragging
          ? "border-primary bg-primary/5"
          : "border-foreground/15 hover:border-foreground/30 hover:bg-foreground/[0.02]",
        loading && "pointer-events-none opacity-60",
      )}
    >
      <input
        ref={inputRef}
        type="file"
        accept=".csv,text/csv"
        className="sr-only"
        onChange={(e) => handleFiles(e.target.files)}
      />
      {loading ? (
        <Loader2 className="size-8 animate-spin text-muted-foreground" />
      ) : fileName ? (
        <FileText className="size-8 text-primary" />
      ) : (
        <Upload className="size-8 text-muted-foreground" />
      )}
      <div className="space-y-1">
        <p className="text-sm font-medium">
          {loading
            ? "Calculating costs…"
            : fileName
              ? fileName
              : "Drop your usage CSV here"}
        </p>
        <p className="text-xs text-muted-foreground">
          {loading
            ? "Crunching your usage data"
            : "or click to browse · .csv up to 10 MB"}
        </p>
      </div>
    </div>
  );
}
