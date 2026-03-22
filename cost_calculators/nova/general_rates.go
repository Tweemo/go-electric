package cost_calculators

import "github.com/tweemo/go-electric/utils"

func NovaGeneralRatesStandardUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Nova", "Basic", "standard")

	totalCost := CalculateGeneralRatesCost(sortedRecords, rate)
	return totalCost
}

func NovaGeneralRatesLowUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Nova", "Basic", "low")

	totalCost := CalculateGeneralRatesCost(sortedRecords, rate)
	return totalCost
}

func CalculateGeneralRatesCost(sortedRecords []utils.DayPower, rate utils.Rate) float64 {
	usage := utils.TotalUsage(sortedRecords)
	levy := utils.GetLevy("Nova")

	cost := (rate.Pwh + levy) * usage
	cost += rate.Daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
