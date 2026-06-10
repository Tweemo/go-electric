// Strips a usage CSV down to an anonymous `timestamp,kwh` file *in the browser*,
// so identifying columns (ICP, meter number, customer name, account number) never
// leave the device. The column heuristics mirror the Go parser in
// backend/utils/usage_data.go so any format the server accepts can be stripped here.

export interface StripResult {
  csv: string; // lean CSV: header "timestamp,kwh" followed by data rows
  rows: number; // number of data rows kept
}

const DELIMITERS = [",", ";", "\t"] as const;

function sniffDelimiter(firstLine: string): string {
  let best = ",";
  let bestCount = -1;
  for (const d of DELIMITERS) {
    const count = firstLine.split(d).length - 1;
    if (count > bestCount) {
      best = d;
      bestCount = count;
    }
  }
  return best;
}

// Minimal RFC4180-ish field splitter: handles quoted fields containing the delimiter.
function parseLine(line: string, delim: string): string[] {
  const out: string[] = [];
  let field = "";
  let inQuotes = false;
  for (let i = 0; i < line.length; i++) {
    const c = line[i];
    if (inQuotes) {
      if (c === '"') {
        if (line[i + 1] === '"') {
          field += '"';
          i++;
        } else {
          inQuotes = false;
        }
      } else {
        field += c;
      }
    } else if (c === '"') {
      inQuotes = true;
    } else if (c === delim) {
      out.push(field);
      field = "";
    } else {
      field += c;
    }
  }
  out.push(field);
  return out;
}

interface ColumnMap {
  dateTimeIdx: number;
  dateIdx: number;
  timeIdx: number;
  usageIdx: number;
  headerRow: boolean;
  hdrDet: boolean;
}

function containsAny(s: string, subs: string[]): boolean {
  return subs.some((sub) => s.includes(sub));
}

function resolveByHeader(header: string[]): ColumnMap | null {
  const cm: ColumnMap = {
    dateTimeIdx: -1,
    dateIdx: -1,
    timeIdx: -1,
    usageIdx: -1,
    headerRow: true,
    hdrDet: false,
  };
  header.forEach((raw, i) => {
    const h = raw.trim().toLowerCase();
    if (cm.usageIdx === -1 && containsAny(h, ["kwh", "usage", "consumption", "value"])) {
      cm.usageIdx = i;
    } else if (cm.dateTimeIdx === -1 && containsAny(h, ["timestamp", "datetime", "date time"])) {
      cm.dateTimeIdx = i;
    } else if (cm.dateIdx === -1 && h.includes("date")) {
      cm.dateIdx = i;
    } else if (cm.timeIdx === -1 && h.includes("time")) {
      cm.timeIdx = i;
    }
  });
  // A standalone date column (no separate time column) is treated as a datetime.
  if (cm.dateTimeIdx === -1 && cm.dateIdx !== -1 && cm.timeIdx === -1) {
    cm.dateTimeIdx = cm.dateIdx;
    cm.dateIdx = -1;
  }
  const hasDate = cm.dateTimeIdx !== -1 || (cm.dateIdx !== -1 && cm.timeIdx !== -1);
  if (hasDate && cm.usageIdx !== -1) return cm;
  return null;
}

function resolveColumns(firstRow: string[]): ColumnMap {
  // Fallback: NZ ICP half-hourly registry file (datetime col 9, usage col 12).
  return (
    resolveByHeader(firstRow) ?? {
      dateTimeIdx: 9,
      dateIdx: -1,
      timeIdx: -1,
      usageIdx: 12,
      headerRow: false,
      hdrDet: true,
    }
  );
}

function csvEscape(v: string): string {
  if (v.includes(",") || v.includes('"') || v.includes("\n")) {
    return `"${v.replace(/"/g, '""')}"`;
  }
  return v;
}

export async function stripCsv(file: File): Promise<StripResult> {
  let text = await file.text();
  if (text.charCodeAt(0) === 0xfeff) text = text.slice(1); // strip UTF-8 BOM

  const lines = text.split(/\r?\n/).filter((l) => l.trim() !== "");
  if (lines.length === 0) throw new Error("The file is empty.");

  const delim = sniffDelimiter(lines[0]);
  const rows = lines.map((l) => parseLine(l, delim));
  const cm = resolveColumns(rows[0]);

  const maxIdx = Math.max(cm.dateTimeIdx, cm.dateIdx, cm.timeIdx, cm.usageIdx);
  const out = ["timestamp,kwh"];
  const start = cm.headerRow ? 1 : 0;

  for (let r = start; r < rows.length; r++) {
    const record = rows[r];
    if (cm.hdrDet && (record.length < 13 || record[0] === "HDR")) continue;
    if (record.length <= maxIdx) continue;

    const dt =
      cm.dateTimeIdx !== -1
        ? record[cm.dateTimeIdx].trim()
        : `${record[cm.dateIdx].trim()} ${record[cm.timeIdx].trim()}`;
    const usage = record[cm.usageIdx].trim();
    if (dt === "" || usage === "") continue;

    out.push(`${csvEscape(dt)},${csvEscape(usage)}`);
  }

  if (out.length === 1) {
    throw new Error(
      "Couldn't find date-time and usage columns in this CSV. Please upload your meter's half-hourly usage export.",
    );
  }

  return { csv: out.join("\n") + "\n", rows: out.length - 1 };
}
