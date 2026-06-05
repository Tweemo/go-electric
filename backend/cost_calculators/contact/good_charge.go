package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func GoodChargeStandardUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodCharge", "standard")

	totalCost := CalculateGoodChargeCost(sortedRecords, rate)
	return totalCost
}

func GoodChargeLowUser(sortedRecords []utils.DayPower) float64 {
	rate := utils.GetRate("Contact", "GoodCharge", "low")

	totalCost := CalculateGoodChargeCost(sortedRecords, rate)
	return totalCost
}

func CalculateGoodChargeCost(sortedRecords []utils.DayPower, rate utils.Rate) float64 {
	weekdayMorningUsage := utils.WeekdayUsage(sortedRecords, 7, 21)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 7, 21)
	morningUsage := weekdayMorningUsage + weekendMorningUsage
	morningCost := morningUsage * rate.Pwh_7am_9pm

	weekdayEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekdayUsage(sortedRecords, 21, 24)
	weekendEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekendUsage(sortedRecords, 21, 24)
	eveningUsage := weekdayEveningUsage + weekendEveningUsage
	eveningCost := eveningUsage * rate.Pwh_9pm_7am

	cost := morningCost + eveningCost + (rate.Daily * 30)

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
