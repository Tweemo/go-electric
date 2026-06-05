package cost_calculators

import (
	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// CalculateGoodChargeCost prices the time-of-use GoodCharge plan. levy is unused.
func CalculateGoodChargeCost(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error) {
	weekdayMorningUsage := utils.WeekdayUsage(records, 7, 21)
	weekendMorningUsage := utils.WeekendUsage(records, 7, 21)
	morningUsage := weekdayMorningUsage + weekendMorningUsage
	morningCost := morningUsage * rate.Pwh7amTo9pm

	weekdayEveningUsage := utils.WeekdayUsage(records, 0, 7) + utils.WeekdayUsage(records, 21, 24)
	weekendEveningUsage := utils.WeekdayUsage(records, 0, 7) + utils.WeekendUsage(records, 21, 24)
	eveningUsage := weekdayEveningUsage + weekendEveningUsage
	eveningCost := eveningUsage * rate.Pwh9pmTo7am

	cost := morningCost + eveningCost + rate.Daily*float64(utils.DayCount(records))

	return utils.RoundFloat(cost, 2)
}
