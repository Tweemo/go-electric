package utils

// DayPower represents power usage data for a day
type DayPower struct {
	date  string
	month string
	day   string
	usage []float64
}

// DayCount returns the number of days represented by the records (one DayPower
// per day), used for daily fixed charges.
func DayCount(records []DayPower) int {
	return len(records)
}

// WeekdayUsage returns the total power usage during weekdays for the given hour
// range across all days in the records.
func WeekdayUsage(records []DayPower, start, end int) float64 {
	usage := 0.0
	for _, record := range records {
		if !isWeekend(record.day) {
			usage += sumSlice(record.usage[start:end])
		}
	}
	return usage
}

// WeekendUsage returns the total power usage during weekends for the given hour
// range across all days in the records.
func WeekendUsage(records []DayPower, start, end int) float64 {
	usage := 0.0
	for _, record := range records {
		if isWeekend(record.day) {
			usage += sumSlice(record.usage[start:end])
		}
	}
	return usage
}

// TotalUsage returns the total power usage across all days in the records.
func TotalUsage(records []DayPower) float64 {
	usage := 0.0
	for _, record := range records {
		usage += sumSlice(record.usage[0:24])
	}
	return usage
}

func isWeekend(day string) bool {
	return day == "Saturday" || day == "Sunday"
}
