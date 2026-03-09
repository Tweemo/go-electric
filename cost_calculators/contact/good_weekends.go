package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func GoodWeekendsStandardUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodWeekends", "standard")

	totalCost := CalculateGoodWeekendsCost(sortedRecords, rate)
	return totalCost
}

func GoodWeekendsLowUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodWeekends", "low")

	totalCost := CalculateGoodWeekendsCost(sortedRecords, rate)
	return totalCost
}

func CalculateGoodWeekendsCost(sortedRecords []utils.DayPower, rate utils.Rate) float64 {
	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 24)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 0, 9)
	weekendEveningUsage := utils.WeekendUsage(sortedRecords, 17, 24)

	usage := weekdayUsage + weekendMorningUsage + weekendEveningUsage

	cost := rate.Pwh * usage
	cost += rate.Daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
