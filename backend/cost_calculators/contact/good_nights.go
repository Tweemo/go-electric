package cost_calculators

import (
	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// CalculateGoodNightsCost prices the GoodNights plan. levy is unused.
func CalculateGoodNightsCost(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error) {
	weekdayUsage := utils.WeekdayUsage(records, 0, 21)
	weekendUsage := utils.WeekendUsage(records, 0, 24)

	usage := weekdayUsage + weekendUsage

	cost := rate.Pwh*usage + rate.Daily*float64(utils.DayCount(records))

	return utils.RoundFloat(cost, 2)
}
