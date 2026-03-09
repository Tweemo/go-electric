package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func GoodNightsStandardUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodNights", "standard")

	totalCost := CalculateGoodNightsCost(sortedRecords, rate)
	return totalCost
}

func GoodNightsLowUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodNights", "low")

	totalCost := CalculateGoodNightsCost(sortedRecords, rate)
	return totalCost
}

func CalculateGoodNightsCost(sortedRecords []utils.DayPower, rate utils.Rate) float64 {
	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 21)
	weekendUsage := utils.WeekendUsage(sortedRecords, 0, 24)

	usage := weekdayUsage + weekendUsage

	cost := rate.Pwh * usage
	cost += rate.Daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
