package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func readCsv(r io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(r)
	csvReader.Comma = ','
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse CSV: %w", err)
	}
	return records, nil
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file %s: %w", filePath, err)
	}
	defer f.Close()
	return readCsv(f)
}

func filterColumns(records [][]string) [][]string {
	var filteredRecords [][]string

	for _, record := range records {
		// Skip header row and empty records
		if len(record) < 13 || record[0] == "HDR" {
			continue
		}

		// Extract only the DateTime columns (9, 10) and float column (12)
		startDateTime := record[9]
		endDateTime := record[10]
		usage := record[12]

		// Validate that usage is a valid float
		if _, err := strconv.ParseFloat(usage, 64); err == nil {
			filteredRecords = append(filteredRecords, []string{startDateTime, endDateTime, usage})
		}
	}

	return filteredRecords
}

// GetUsageData reads and filters usage data from a CSV file on disk.
func GetUsageData(filepath string) ([][]string, error) {
	records, err := readCsvFile(filepath)
	if err != nil {
		return nil, err
	}
	return filterColumns(records), nil
}

// ParseUsageData reads and filters usage data from an arbitrary reader (e.g. an
// uploaded multipart file), avoiding a round-trip through disk.
func ParseUsageData(r io.Reader) ([][]string, error) {
	records, err := readCsv(r)
	if err != nil {
		return nil, err
	}
	return filterColumns(records), nil
}

// CalculateDayPower calculates the power usage for each day
func CalculateDayPower(usageData [][]string) ([]DayPower, error) {
	var sortedRecords []DayPower
	i := 0

	for _, record := range usageData {
		startTime, _ := time.Parse("02/01/2006 15:04:05", record[0])
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
