package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func CalculateGoodNightsStandardUserCost(sortedRecords []utils.DayPower) float64 {
	ceGnSuPwhCharge := utils.MustFloat64Env("CE_GN_SU_PWH_CHARGE")
	ceGnSuDailyCharge := utils.MustFloat64Env("CE_GN_SU_DAILY_CHARGE")

	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 21)
	weekendUsage := utils.WeekendUsage(sortedRecords, 0, 24)

	usage := weekdayUsage + weekendUsage

	cost := ceGnSuPwhCharge * usage
	cost += ceGnSuDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}

func CalculateGoodNightsLowUserCost(sortedRecords []utils.DayPower) float64 {
	ceGnLuPwhCharge := utils.MustFloat64Env("CE_GN_LU_PWH_CHARGE")
	ceGnLuDailyCharge := utils.MustFloat64Env("CE_GN_LU_DAILY_CHARGE")

	weekdayUsage := utils.WeekdayUsage(sortedRecords, 0, 21)
	weekendUsage := utils.WeekendUsage(sortedRecords, 0, 24)

	usage := weekdayUsage + weekendUsage

	cost := ceGnLuPwhCharge * usage
	cost += ceGnLuDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
