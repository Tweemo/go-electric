package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func SimpleRatesStandardUser(sortedRecords []utils.DayPower) float64 {
	standard := utils.GetRate("Contact", "GoodNights", "standard")
	standardRateMap := standard.(map[string]interface{})

	usage := utils.TotalUsage(sortedRecords)
	pwh, _ := utils.GetFloat(standardRateMap["pwh"])
	daily, _ := utils.GetFloat(standardRateMap["daily"])

	totalCost := CalculateSimpleRatesCost(usage, pwh, daily)
	return totalCost
}

func SimpleRatesLowUser(sortedRecords []utils.DayPower) float64 {
	low := utils.GetRate("Contact", "GoodNights", "low")
	lowRateMap := low.(map[string]interface{})

	usage := utils.TotalUsage(sortedRecords)
	pwh, _ := utils.GetFloat(lowRateMap["pwh"])
	daily, _ := utils.GetFloat(lowRateMap["daily"])

	totalCost := CalculateSimpleRatesCost(usage, pwh, daily)
	return totalCost
}

func CalculateSimpleRatesCost(usage float64, pwh float64, daily float64) float64 {
	levy := utils.GetLevy("Contact")

	cost := (pwh + levy) * usage
	cost += daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
