package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func GoodWeekendsStandardUser(sortedRecords []utils.DayPower) float64 {
	standard := utils.GetRate("Contact", "GoodWeekends", "standard")
	standardRateMap := standard.(map[string]interface{})

	pwh, _ := utils.GetFloat(standardRateMap["pwh"])
	daily, _ := utils.GetFloat(standardRateMap["daily"])

	totalCost := CalculateGoodWeekendsCost(sortedRecords, pwh, daily)
	return totalCost
}

func GoodWeekendsLowUser(sortedRecords []utils.DayPower) float64 {
	low := utils.GetRate("Contact", "GoodWeekends", "low")
	lowRateMap := low.(map[string]interface{})

	pwh, _ := utils.GetFloat(lowRateMap["pwh"])
	daily, _ := utils.GetFloat(lowRateMap["daily"])

	totalCost := CalculateGoodWeekendsCost(sortedRecords, pwh, daily)
	return totalCost
}

func CalculateGoodWeekendsCost(sortedRecords []utils.DayPower, pwh float64, daily float64) float64 {
	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 24)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 0, 9)
	weekendEveningUsage := utils.WeekendUsage(sortedRecords, 17, 24)

	usage := weekdayUsage + weekendMorningUsage + weekendEveningUsage

	cost := pwh * usage
	cost += daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
