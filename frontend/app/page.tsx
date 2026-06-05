"use client";

import * as React from "react";
import { TriangleAlert } from "lucide-react";

import { UploadDropzone } from "@/components/UploadDropzone";
import { ResultsView } from "@/components/ResultsView";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { postCosts } from "@/lib/api";
import type { CostsResponse } from "@/lib/plans";

export default function Home() {
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [fileName, setFileName] = React.useState<string | null>(null);
  const [data, setData] = React.useState<CostsResponse | null>(null);

  async function handleFile(file: File) {
    setLoading(true);
    setError(null);
    setFileName(file.name);
    try {
      const result = await postCosts(file);
      setData(result);
    } catch (err) {
      setData(null);
      setError(err instanceof Error ? err.message : "Something went wrong.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto w-full max-w-3xl px-6 py-16">
      <header className="mb-8 space-y-2">
        <h1 className="text-2xl font-semibold tracking-tight">
          Power plan cost comparison
        </h1>
        <p className="text-sm text-muted-foreground">
          Upload your electricity usage CSV to see what each plan would cost.
        </p>
      </header>

      <UploadDropzone
        onFile={handleFile}
        loading={loading}
        fileName={fileName}
      />

      {error && (
        <Alert variant="destructive" className="mt-6">
          <TriangleAlert />
          <AlertTitle>Couldn&apos;t calculate costs</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {data && (
        <div className="mt-10">
          <ResultsView data={data} />
        </div>
      )}
    </main>
  );
}
