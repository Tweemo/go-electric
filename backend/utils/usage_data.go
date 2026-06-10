package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// canonicalLayout is the datetime format CalculateDayPower expects. Every parsed
// row is re-emitted in this layout so downstream aggregation stays format-agnostic.
const canonicalLayout = "02/01/2006 15:04:05"

// ParseResult is the outcome of reading a usage CSV: the normalized rows plus
// counts so callers can tell the user how much of their file was usable.
type ParseResult struct {
	Records     [][]string // each: {canonicalDateTime, "", usage}
	RowsParsed  int
	RowsSkipped int
}

// timestampLayouts are tried in order against each row's datetime cell, covering
// the common NZ retailer/meter exports (DD/MM/YYYY) and ISO-8601 variants.
var timestampLayouts = []string{
	canonicalLayout,
	"02/01/2006 15:04",
	"02/01/2006",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	time.RFC3339,
}

func parseTimestamp(s string) (time.Time, error) {
	for _, layout := range timestampLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized timestamp %q", s)
}

// readCsv reads an entire CSV into memory, tolerating a UTF-8 BOM, ragged rows,
// leading spaces, and comma/semicolon/tab delimiters.
func readCsv(r io.Reader) ([][]string, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read CSV: %w", err)
	}
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf")) // strip UTF-8 BOM

	csvReader := csv.NewReader(bytes.NewReader(raw))
	csvReader.Comma = sniffDelimiter(raw)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse CSV: %w", err)
	}
	return records, nil
}

// sniffDelimiter inspects the first non-empty line and picks whichever of
// comma/semicolon/tab appears most often, defaulting to comma.
func sniffDelimiter(raw []byte) rune {
	for line := range strings.SplitSeq(string(raw), "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		delim, best := ',', strings.Count(line, ",")
		if n := strings.Count(line, ";"); n > best {
			delim, best = ';', n
		}
		if n := strings.Count(line, "\t"); n > best {
			delim = '\t'
		}
		return delim
	}
	return ','
}

// columnMap describes where the datetime and usage values live in a CSV.
type columnMap struct {
	dateTimeIdx int // combined datetime column, or -1
	dateIdx     int // separate date column, or -1
	timeIdx     int // separate time column, or -1
	usageIdx    int
	headerRow   bool // first row is a header to skip
	hdrDet      bool // NZ ICP "HDR/DET" registry convention
}

// resolveColumns decides how to read a file: by header names when the first row
// looks like a header, otherwise via the HDR/DET registry convention.
func resolveColumns(records [][]string) (columnMap, error) {
	if len(records) == 0 {
		return columnMap{}, errors.New("CSV is empty")
	}
	if cm, ok := resolveByHeader(records[0]); ok {
		return cm, nil
	}
	// Fallback: NZ ICP half-hourly registry file (datetime col 9, usage col 12).
	return columnMap{dateTimeIdx: 9, dateIdx: -1, timeIdx: -1, usageIdx: 12, hdrDet: true}, nil
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// resolveByHeader matches columns by case-insensitive header keywords. It returns
// ok=false when no usage/datetime columns can be identified (e.g. a HDR/DET file).
func resolveByHeader(header []string) (columnMap, bool) {
	cm := columnMap{dateTimeIdx: -1, dateIdx: -1, timeIdx: -1, usageIdx: -1, headerRow: true}
	for i, raw := range header {
		h := strings.ToLower(strings.TrimSpace(raw))
		switch {
		case cm.usageIdx == -1 && containsAny(h, "kwh", "usage", "consumption", "value"):
			cm.usageIdx = i
		case cm.dateTimeIdx == -1 && containsAny(h, "timestamp", "datetime", "date time"):
			cm.dateTimeIdx = i
		case cm.dateIdx == -1 && strings.Contains(h, "date"):
			cm.dateIdx = i
		case cm.timeIdx == -1 && strings.Contains(h, "time"):
			cm.timeIdx = i
		}
	}
	// A standalone date column (no separate time column) is treated as a datetime.
	if cm.dateTimeIdx == -1 && cm.dateIdx != -1 && cm.timeIdx == -1 {
		cm.dateTimeIdx, cm.dateIdx = cm.dateIdx, -1
	}
	hasDate := cm.dateTimeIdx != -1 || (cm.dateIdx != -1 && cm.timeIdx != -1)
	if hasDate && cm.usageIdx != -1 {
		return cm, true
	}
	return columnMap{}, false
}

// extractRecords normalizes data rows into {canonicalDateTime, "", usage}, counting
// rows it has to skip and erroring only when nothing usable remains.
func extractRecords(records [][]string, cm columnMap) (ParseResult, error) {
	maxIdx := cm.usageIdx
	for _, idx := range []int{cm.dateTimeIdx, cm.dateIdx, cm.timeIdx} {
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	start := 0
	if cm.headerRow {
		start = 1
	}

	var out [][]string
	skipped := 0
	for r := start; r < len(records); r++ {
		record := records[r]

		// Structural rows in the registry format are not data rows; don't count them.
		if cm.hdrDet && (len(record) < 13 || (len(record) > 0 && record[0] == "HDR")) {
			continue
		}
		if len(record) <= maxIdx {
			skipped++
			continue
		}

		var dtStr string
		if cm.dateTimeIdx != -1 {
			dtStr = strings.TrimSpace(record[cm.dateTimeIdx])
		} else {
			dtStr = strings.TrimSpace(record[cm.dateIdx]) + " " + strings.TrimSpace(record[cm.timeIdx])
		}
		ts, err := parseTimestamp(dtStr)
		if err != nil {
			skipped++
			continue
		}

		usage := strings.TrimSpace(record[cm.usageIdx])
		if _, err := strconv.ParseFloat(usage, 64); err != nil {
			skipped++
			continue
		}

		out = append(out, []string{ts.Format(canonicalLayout), "", usage})
	}

	if len(out) == 0 {
		return ParseResult{}, errors.New("no valid usage rows found")
	}
	return ParseResult{Records: out, RowsParsed: len(out), RowsSkipped: skipped}, nil
}

// GetUsageData reads and normalizes usage data from a CSV file on disk.
func GetUsageData(filepath string) (ParseResult, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return ParseResult{}, fmt.Errorf("unable to read input file %s: %w", filepath, err)
	}
	defer f.Close()
	return ParseUsageData(f)
}

// ParseUsageData reads and normalizes usage data from an arbitrary reader (e.g. an
// uploaded multipart file), avoiding a round-trip through disk.
func ParseUsageData(r io.Reader) (ParseResult, error) {
	records, err := readCsv(r)
	if err != nil {
		return ParseResult{}, err
	}
	cm, err := resolveColumns(records)
	if err != nil {
		return ParseResult{}, err
	}
	return extractRecords(records, cm)
}

// CalculateDayPower calculates the power usage for each day
func CalculateDayPower(usageData [][]string) ([]DayPower, error) {
	var sortedRecords []DayPower
	i := 0

	for _, record := range usageData {
		startTime, _ := time.Parse(canonicalLayout, record[0])
		month := startTime.Format("01/2006")
		date := startTime.Format("02/01/2006")
		time := startTime.Format("15:04:05")
		day := startTime.Format("Monday")
		usage, _ := strconv.ParseFloat(record[2], 64)
		hour := startTime.Hour()

		// First record in the data
		if len(sortedRecords) == 0 {
			dayUsage := make([]float64, 24)
			dayUsage[hour] += usage
			sortedRecords = append(sortedRecords, DayPower{date: date, month: month, day: day, usage: dayUsage})
		} else if i < len(sortedRecords) && sortedRecords[i].date == date {
			// Accumulate into the hour bucket (sum both half-hours into the hour)
			if hour >= 0 && hour < 24 {
				sortedRecords[i].usage[hour] += usage
			}
			if time == "23:30:00" {
				i++
			}
		} else {
			dayUsage := make([]float64, 24)
			dayUsage[hour] += usage
			sortedRecords = append(sortedRecords, DayPower{date: date, month: month, day: day, usage: dayUsage})
		}
	}

	return sortedRecords, nil
}

func sumSlice(numbers []float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += number
	}
	return total
}
