import type { CostsResponse } from "@/lib/plans";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// Posts the CSV to the Go /costs endpoint. The form field MUST be "file" to
// match c.FormFile("file") on the server.
export async function postCosts(file: File): Promise<CostsResponse> {
  const form = new FormData();
  form.append("file", file);

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
