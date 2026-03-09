package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func SimpleRatesStandardUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodNights", "standard")

	totalCost := CalculateSimpleRatesCost(sortedRecords, rate)
	return totalCost
}

func SimpleRatesLowUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodNights", "low")

	totalCost := CalculateSimpleRatesCost(sortedRecords, rate)
	return totalCost
}

func CalculateSimpleRatesCost(sortedRecords []utils.DayPower, rate utils.Rate) float64 {
	levy := utils.GetLevy("Contact")
	usage := utils.TotalUsage(sortedRecords)

	cost := (rate.Pwh + levy) * usage
	cost += rate.Daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
