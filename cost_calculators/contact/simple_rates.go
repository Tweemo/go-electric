package cost_calculators

import (
	"github.com/tweemo/go-electric/utils"
)

func CalculateSimpleRatesStandardUserCost(usage float64) float64 {
	ceSrSuPwhCharge := utils.MustFloat64Env("CE_SR_SU_PWH_CHARGE")
	ceSrSuDailyCharge := utils.MustFloat64Env("CE_SR_SU_DAILY_CHARGE")
	ceLevy := utils.MustFloat64Env("CE_LEVY")

	cost := (ceSrSuPwhCharge + ceLevy) * usage
	cost += ceSrSuDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}

func CalculateSimpleRatesLowUserCost(usage float64) float64 {
	ceSrLuPwhCharge := utils.MustFloat64Env("CE_SR_LU_PWH_CHARGE")
	ceSrLuDailyCharge := utils.MustFloat64Env("CE_SR_LU_DAILY_CHARGE")
	ceLevy := utils.MustFloat64Env("CE_LEVY")

	cost := (ceSrLuPwhCharge + ceLevy) * usage
	cost += ceSrLuDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
