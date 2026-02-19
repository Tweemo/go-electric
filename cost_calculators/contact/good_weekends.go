package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func CalculateGoodWeekendsStandardUserCost(sortedRecords []utils.DayPower) float64 {
	ceGwSuPwhCharge := utils.MustFloat64Env("CE_GW_SU_PWH_CHARGE")
	ceGwSuDailyCharge := utils.MustFloat64Env("CE_GW_SU_DAILY_CHARGE")

	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 24)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 0, 9)
	weekendEveningUsage := utils.WeekendUsage(sortedRecords, 17, 24)

	usage := weekdayUsage + weekendMorningUsage + weekendEveningUsage

	cost := ceGwSuPwhCharge * usage
	cost += ceGwSuDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}

func CalculateGoodWeekendsLowUserCost(sortedRecords []utils.DayPower) float64 {
	ceGwLuPwhCharge := utils.MustFloat64Env("CE_GW_LU_PWH_CHARGE")
	ceGwLuDailyCharge := utils.MustFloat64Env("CE_GW_LU_DAILY_CHARGE")

	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 24)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 0, 9)
	weekendEveningUsage := utils.WeekendUsage(sortedRecords, 17, 24)

	usage := weekdayUsage + weekendMorningUsage + weekendEveningUsage

	cost := ceGwLuPwhCharge * usage
	cost += ceGwLuDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
