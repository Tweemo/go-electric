import type { CostsResponse } from "@/lib/plans";
import { stripCsv } from "@/lib/stripCsv";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// Posts the CSV to the Go /costs endpoint. The form field MUST be "file" to
// match c.FormFile("file") on the server.
//
// The file is stripped to an anonymous `timestamp,kwh` CSV in the browser first,
// so identifiers (ICP, meter number, name, account) never leave the device — only
// the date-time and usage figures needed to compute costs are uploaded.
export async function postCosts(file: File): Promise<CostsResponse> {
  const { csv } = await stripCsv(file);
  const lean = new File([csv], "usage.csv", { type: "text/csv" });

  const form = new FormData();
  form.append("file", lean);

  let res: Response;
  try {
    res = await fetch(`${API_URL}/costs`, { method: "POST", body: form });
  } catch {
    throw new Error(
      `Could not reach the API at ${API_URL}. Is the Go server running?`,
    );
  }

  let body: unknown = null;
  try {
    body = await res.json();
  } catch {
    // Non-JSON response; fall through to status handling below.
  }

  if (!res.ok) {
    const message =
      body && typeof body === "object" && "error" in body
        ? String((body as { error: unknown }).error)
        : `Request failed (${res.status})`;
    throw new Error(message);
  }

  return body as CostsResponse;
}
