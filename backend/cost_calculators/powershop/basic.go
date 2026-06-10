package cost_calculators

import (
	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// CalculateBasicCost prices Powershop's time-of-use Basic plan. Peak rate applies
// 7am-9pm and off-peak 9pm-7am, with no weekday/weekend distinction. levy is unused
// (Powershop charges no levy).
func CalculateBasicCost(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error) {
	peakUsage := utils.WeekdayUsage(records, 7, 21) + utils.WeekendUsage(records, 7, 21)
	peakCost := peakUsage * rate.Pwh7amTo9pm

	offPeakUsage := utils.WeekdayUsage(records, 0, 7) + utils.WeekdayUsage(records, 21, 24) +
		utils.WeekendUsage(records, 0, 7) + utils.WeekendUsage(records, 21, 24)
	offPeakCost := offPeakUsage * rate.Pwh9pmTo7am

	cost := peakCost + offPeakCost + rate.Daily*float64(utils.DayCount(records))

	return utils.RoundFloat(cost, 2)
}
