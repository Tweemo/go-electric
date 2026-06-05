package cost_calculators

import (
	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// CalculateGoodWeekendsCost prices the GoodWeekends plan. levy is unused.
func CalculateGoodWeekendsCost(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error) {
	weekdayUsage := utils.WeekdayUsage(records, 0, 24)
	weekendMorningUsage := utils.WeekendUsage(records, 0, 9)
	weekendEveningUsage := utils.WeekendUsage(records, 17, 24)

	usage := weekdayUsage + weekendMorningUsage + weekendEveningUsage

	cost := rate.Pwh*usage + rate.Daily*float64(utils.DayCount(records))

	return utils.RoundFloat(cost, 2)
}
