package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func GoodChargeStandardUser(sortedRecords []utils.DayPower) float64 {
	standard := utils.GetRate("Contact", "GoodCharge", "standard")
	standardRateMap := standard.(map[string]interface{})

	pwh7am9pm, _ := utils.GetFloat(standardRateMap["pwh_7am_9pm"])
	pwh9pm7am, _ := utils.GetFloat(standardRateMap["pwh_9pm_7am"])
	daily, _ := utils.GetFloat(standardRateMap["daily"])

	totalCost := CalculateGoodChargeCost(sortedRecords, pwh7am9pm, pwh9pm7am, daily)
	return totalCost
}

func GoodChargeLowUser(sortedRecords []utils.DayPower) float64 {
	low := utils.GetRate("Contact", "GoodCharge", "standard")
	lowRateMap := low.(map[string]interface{})

	pwh7am9pm, _ := utils.GetFloat(lowRateMap["pwh_7am_9pm"])
	pwh9pm7am, _ := utils.GetFloat(lowRateMap["pwh_9pm_7am"])
	daily, _ := utils.GetFloat(lowRateMap["daily"])

	totalCost := CalculateGoodChargeCost(sortedRecords, pwh7am9pm, pwh9pm7am, daily)
	return totalCost
}

func CalculateGoodChargeCost(sortedRecords []utils.DayPower, pwh7am9pm float64, pwh9pm7am float64, daily float64) float64 {
	weekdayMorningUsage := utils.WeekdayUsage(sortedRecords, 7, 21)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 7, 21)
	morningUsage := weekdayMorningUsage + weekendMorningUsage
	morningCost := morningUsage * pwh7am9pm

	weekdayEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekdayUsage(sortedRecords, 21, 24)
	weekendEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekendUsage(sortedRecords, 21, 24)
	eveningUsage := weekdayEveningUsage + weekendEveningUsage
	eveningCost := eveningUsage * pwh9pm7am

	cost := morningCost + eveningCost + (daily * 30)

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
