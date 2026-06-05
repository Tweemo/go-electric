package cost_calculators

import (
	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// CalculateGeneralRatesCost prices Nova's general rates plan, applying the levy.
func CalculateGeneralRatesCost(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error) {
	usage := utils.TotalUsage(records)

	cost := (rate.Pwh+levy)*usage + rate.Daily*float64(utils.DayCount(records))

	return utils.RoundFloat(cost, 2)
}
