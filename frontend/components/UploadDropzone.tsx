"use client";

import * as React from "react";
import { Upload, Download, Check, RefreshCw } from "lucide-react";

import { cn } from "@/lib/utils";

interface UploadDropzoneProps {
  onFile: (file: File) => void;
  onSample: () => void;
  onReset: () => void;
  loading: boolean;
  fileName: string | null;
  /** When true, show the compact "file selected" bar instead of the drop area. */
  selected: boolean;
}

export function UploadDropzone({
  onFile,
  onSample,
  onReset,
  loading,
  fileName,
  selected,
}: UploadDropzoneProps) {
  const inputRef = React.useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = React.useState(false);

  function handleFiles(files: FileList | null) {
    const file = files?.[0];
    if (file) onFile(file);
  }

  const fileInput = (
    <input
      ref={inputRef}
      type="file"
      accept=".csv,text/csv"
      className="sr-only"
      onChange={(e) => handleFiles(e.target.files)}
    />
  );

  // Compact bar once a file has produced results.
  if (selected) {
    return (
      <div className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-primary/40 bg-accent px-5 py-3.5">
        {fileInput}
        <div className="flex items-center gap-2.5">
          <div className="flex size-7 items-center justify-center rounded-lg bg-primary">
            <Check size={14} className="text-primary-foreground" />
          </div>
          <span className="text-sm font-semibold text-accent-foreground">
            {fileName}
          </span>
        </div>
        <button
          type="button"
          onClick={onReset}
          className="flex items-center gap-1.5 rounded-lg border-[1.5px] border-primary/40 px-3.5 py-1.5 text-[13px] font-semibold text-accent-foreground transition-colors hover:bg-card/50"
        >
          <RefreshCw size={12} /> Upload different file
        </button>
      </div>
    );
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
        "flex w-full cursor-pointer flex-col items-center justify-center gap-4 rounded-2xl border-2 border-dashed px-8 py-14 text-center transition-colors",
        dragging
          ? "border-primary bg-accent"
          : "border-border bg-card hover:border-primary/50",
        loading && "pointer-events-none opacity-70",
      )}
    >
      {fileInput}
      {loading ? (
        <div className="flex flex-col items-center gap-4">
          <div className="flex size-14 items-center justify-center rounded-full bg-accent">
            <div className="size-6 animate-spin rounded-full border-[3px] border-primary border-t-transparent" />
          </div>
          <div className="text-base font-semibold text-foreground">
            Crunching your numbers…
          </div>
        </div>
      ) : (
        <>
          <div
            className={cn(
              "flex size-16 items-center justify-center rounded-2xl transition-colors",
              dragging ? "bg-primary" : "bg-secondary",
            )}
          >
            <Upload
              size={28}
              className={dragging ? "text-primary-foreground" : "text-primary"}
            />
          </div>
          <div>
            <div className="mb-1.5 text-xl font-bold text-heading">
              Drop your usage CSV here
            </div>
            <div className="text-[15px] text-subtle">
              or click to browse — .csv up to 10 MB
            </div>
          </div>
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              onSample();
            }}
            className="flex items-center gap-1.5 text-[13px] font-medium text-primary hover:underline"
          >
            <Download size={13} /> Download a sample CSV to try it out
          </button>
        </>
      )}
    </div>
  );
}
