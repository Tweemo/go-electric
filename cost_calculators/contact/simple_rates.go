package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func SimpleRatesStandardUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodNights", "standard")
	usage := utils.TotalUsage(sortedRecords)

	totalCost := CalculateSimpleRatesCost(usage, rate)
	return totalCost
}

func SimpleRatesLowUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodNights", "low")
	usage := utils.TotalUsage(sortedRecords)

	totalCost := CalculateSimpleRatesCost(usage, rate)
	return totalCost
}

func CalculateSimpleRatesCost(usage float64, rate utils.Rate) float64 {
	levy := utils.GetLevy("Contact")

	cost := (rate.Pwh + levy) * usage
	cost += rate.Daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
