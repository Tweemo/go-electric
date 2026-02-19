package utils

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func createMonthMap(sortedRecords []DayPower) map[string][]DayPower {
	monthMap := make(map[string][]DayPower)
	for _, record := range sortedRecords {
		monthMap[record.month] = append(monthMap[record.month], record)
	}

	return monthMap
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func filterColumns(records [][]string) [][]string {
	var filteredRecords [][]string

	for _, record := range records {
		// Skip header row and empty records
		if len(record) < 13 || record[0] == "HDR" {
			continue
		}

		// Extract only the DateTime columns (9, 10) and float column (12)
		if len(record) >= 13 {
			startDateTime := record[9]
			endDateTime := record[10]
			usage := record[12]

			// Validate that usage is a valid float
			if _, err := strconv.ParseFloat(usage, 64); err == nil {
				filteredRecords = append(filteredRecords, []string{startDateTime, endDateTime, usage})
			}
		}
	}

	return filteredRecords
}

// getUsageData returns the usage data from some source
func GetUsageData() [][]string {
	dirname, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get user home directory", err)
	}

	usageFileName := os.Getenv("USAGE_FILE_NAME")
	if usageFileName == "" {
		log.Fatal("USAGE_FILE_NAME is not set in the .env file")
	}

	records := readCsvFile(filepath.Join(dirname, "data", usageFileName))

	// Filter to keep only DateTime columns and float column
	filteredRecords := filterColumns(records)

	return filteredRecords
}

// calculateDayPower calculates the power usage for each day
func CalculateDayPower(usageData [][]string) []DayPower {
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

	return sortedRecords
}

func sumSlice(numbers []float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += number
	}
	return total
}
