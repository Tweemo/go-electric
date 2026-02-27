package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func GoodNightsStandardUser(sortedRecords []utils.DayPower) float64 {
	standard := utils.GetRate("Contact", "GoodNights", "standard")
	standardRateMap := standard.(map[string]interface{})

	pwh, _ := utils.GetFloat(standardRateMap["pwh"])
	daily, _ := utils.GetFloat(standardRateMap["daily"])

	totalCost := CalculateGoodNightsCost(sortedRecords, pwh, daily)
	return totalCost
}

func GoodNightsLowUser(sortedRecords []utils.DayPower) float64 {
	low := utils.GetRate("Contact", "GoodNights", "low")
	lowRateMap := low.(map[string]interface{})

	pwh, _ := utils.GetFloat(lowRateMap["pwh"])
	daily, _ := utils.GetFloat(lowRateMap["daily"])

	totalCost := CalculateGoodNightsCost(sortedRecords, pwh, daily)
	return totalCost
}

func CalculateGoodNightsCost(sortedRecords []utils.DayPower, pwh float64, daily float64) float64 {
	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 21)
	weekendUsage := utils.WeekendUsage(sortedRecords, 0, 24)

	usage := weekdayUsage + weekendUsage

	cost := pwh * usage
	cost += daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
