package cost_calculatr

import "github.com/tweemo/go-electric/utils"

func CalculateNovaStandardUserCost(sortedRecords []utils.DayPower) float64 {
	nvSuPwhCharge := utils.MustFloat64Env("NV_SU_PWH_CHARGE")
	nvDailyCharge := utils.MustFloat64Env("NV_SU_DAILY_CHARGE")
	nvLevy := utils.MustFloat64Env("NV_LEVY")

	usage := utils.TotalUsage(sortedRecords)

	cost := (nvSuPwhCharge + nvLevy) * usage
	cost += nvDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}

func CalculateNovaLowUserCost(sortedRecords []utils.DayPower) float64 {
	nvSuPwhCharge := utils.MustFloat64Env("NV_LU_PWH_CHARGE")
	nvDailyCharge := utils.MustFloat64Env("NV_LU_DAILY_CHARGE")
	nvLevy := utils.MustFloat64Env("NV_LEVY")

	usage := utils.TotalUsage(sortedRecords)

	cost := (nvSuPwhCharge + nvLevy) * usage
	cost += nvDailyCharge * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
