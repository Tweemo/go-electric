package utils

// DayPower represents power usage data for a day
type DayPower struct {
	date  string
	month string
	day   string
	usage []float64
}

// WeekdayUsage returns the total power usage during weekdays
func WeekdayUsage(records []DayPower, start, end int) float64 {
	usage := 0.0

	monthMap := createMonthMap(records)
	monthRecords := monthMap["08/2025"]
	for _, record := range monthRecords {
		if !isWeekend(record.day) {
			usage += sumSlice(record.usage[start:end])
		}
	}

	return usage
}

// WeekendUsage returns the total power usage during weekends
func WeekendUsage(records []DayPower, start, end int) float64 {
	usage := 0.0

	monthMap := createMonthMap(records)
	monthRecords := monthMap["08/2025"]
	for _, record := range monthRecords {
		if isWeekend(record.day) {
			usage += sumSlice(record.usage[start:end])
		}
	}

	return usage
}

// TotalUsage returns the total power usage across all days
func TotalUsage(records []DayPower) float64 {
	usage := 0.0

	monthMap := createMonthMap(records)
	monthRecords := monthMap["08/2025"]
	for _, record := range monthRecords {
		usage += sumSlice(record.usage[0:24])
	}

	return usage
}

func isWeekend(day string) bool {
	if day == "Saturday" || day == "Sunday" {
		return true
	}

	return false
}
