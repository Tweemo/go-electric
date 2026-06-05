"use client";

import * as React from "react";
import { TriangleAlert } from "lucide-react";

import { SiteNav } from "@/components/SiteNav";
import { Hero } from "@/components/Hero";
import { UploadDropzone } from "@/components/UploadDropzone";
import { ResultsView } from "@/components/ResultsView";
import { HowItWorks } from "@/components/HowItWorks";
import { TrustStrip } from "@/components/TrustStrip";
import { CsvFormatGuide } from "@/components/CsvFormatGuide";
import { SiteFooter } from "@/components/SiteFooter";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { postCosts } from "@/lib/api";
import type { CostsResponse } from "@/lib/plans";

export default function Home() {
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [fileName, setFileName] = React.useState<string | null>(null);
  const [data, setData] = React.useState<CostsResponse | null>(null);

  const heroInputRef = React.useRef<HTMLInputElement>(null);
  const resultsRef = React.useRef<HTMLElement>(null);

  async function handleFile(file: File) {
    setLoading(true);
    setError(null);
    setFileName(file.name);
    try {
      const result = await postCosts(file);
      setData(result);
      setTimeout(
        () => resultsRef.current?.scrollIntoView({ behavior: "smooth", block: "start" }),
        120,
      );
    } catch (err) {
      setData(null);
      setError(err instanceof Error ? err.message : "Something went wrong.");
    } finally {
      setLoading(false);
    }
  }

  function triggerUpload() {
    heroInputRef.current?.click();
  }

  function downloadSample() {
    const a = document.createElement("a");
    a.href = "/sample-usage.csv";
    a.download = "sample-usage.csv";
    a.click();
  }

  function reset() {
    setData(null);
    setError(null);
    setFileName(null);
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="flex min-h-full flex-col">
      <SiteNav onSample={downloadSample} />
      <Hero onUpload={triggerUpload} onSample={downloadSample} />

      {/* Hidden input for the nav/hero "Upload" buttons. */}
      <input
        ref={heroInputRef}
        type="file"
        accept=".csv,text/csv"
        className="sr-only"
        onChange={(e) => {
          const f = e.target.files?.[0];
          if (f) handleFile(f);
        }}
      />

      <section className="mx-auto w-full max-w-[1080px] px-7 pb-20">
        <UploadDropzone
          onFile={handleFile}
          onSample={downloadSample}
          onReset={reset}
          loading={loading}
          fileName={fileName}
          selected={Boolean(data)}
        />

        {error && (
          <Alert variant="destructive" className="mt-4">
            <TriangleAlert />
            <AlertTitle>Couldn&apos;t calculate costs</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
      </section>

      {data && <ResultsView ref={resultsRef} data={data} />}

      <HowItWorks />
      <TrustStrip />
      <CsvFormatGuide onSample={downloadSample} />
      <SiteFooter />
    </div>
  );
}
