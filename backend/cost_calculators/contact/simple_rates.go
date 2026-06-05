package cost_calculators

import (
	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// CalculateSimpleRatesCost prices the SimpleRates plan, applying the company levy.
func CalculateSimpleRatesCost(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error) {
	usage := utils.TotalUsage(records)

	cost := (rate.Pwh+levy)*usage + rate.Daily*float64(utils.DayCount(records))

	return utils.RoundFloat(cost, 2)
}
