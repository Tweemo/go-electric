package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func CalculateGoodChargeStandardUserCost(sortedRecords []utils.DayPower) float64 {
	ceGcSuPwhCharge7AM9PM := utils.MustFloat64Env("CE_GC_SU_PWH_7AM_9PM_CHARGE")
	ceGcSuPwhCharge9PM7AM := utils.MustFloat64Env("CE_GC_SU_PWH_9PM_7AM_CHARGE")
	ceGcSuDailyCharge := utils.MustFloat64Env("CE_GC_SU_DAILY_CHARGE")

	weekdayMorningUsage := utils.WeekdayUsage(sortedRecords, 7, 21)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 7, 21)
	morningUsage := weekdayMorningUsage + weekendMorningUsage
	morningCost := morningUsage * ceGcSuPwhCharge7AM9PM

	weekdayEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekdayUsage(sortedRecords, 21, 24)
	weekendEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekendUsage(sortedRecords, 21, 24)
	eveningUsage := weekdayEveningUsage + weekendEveningUsage
	eveningCost := eveningUsage * ceGcSuPwhCharge9PM7AM

	cost := morningCost + eveningCost + (ceGcSuDailyCharge * 30)

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}

func CalculateGoodChargeLowUserCost(sortedRecords []utils.DayPower) float64 {
	ceGcLuPwhCharge7AM9PM := utils.MustFloat64Env("CE_GC_LU_PWH_7AM_9PM_CHARGE")
	ceGcLuPwhCharge9PM7AM := utils.MustFloat64Env("CE_GC_LU_PWH_9PM_7AM_CHARGE")
	ceGcLuDailyCharge := utils.MustFloat64Env("CE_GC_LU_DAILY_CHARGE")

	weekdayMorningUsage := utils.WeekdayUsage(sortedRecords, 7, 21)
	weekendMorningUsage := utils.WeekendUsage(sortedRecords, 7, 21)
	morningUsage := weekdayMorningUsage + weekendMorningUsage
	morningCost := morningUsage * ceGcLuPwhCharge7AM9PM

	weekdayEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekdayUsage(sortedRecords, 21, 24)
	weekendEveningUsage := utils.WeekdayUsage(sortedRecords, 0, 7) + utils.WeekendUsage(sortedRecords, 21, 24)
	eveningUsage := weekdayEveningUsage + weekendEveningUsage
	eveningCost := eveningUsage * ceGcLuPwhCharge9PM7AM

	cost := morningCost + eveningCost + (ceGcLuDailyCharge * 30)

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
