package cost_calculators

import (
	costs "github.com/tweemo/go-electric/cost_calculators"
	"github.com/tweemo/go-electric/utils"
)

func NovaCosts(sortedRecords []utils.DayPower) [1]costs.Cost {
	var cost [1]costs.Cost

	// General Rates
	GeneralRatesStandard := NovaGeneralRatesStandardUser(sortedRecords)
	GeneralRatesLow := NovaGeneralRatesLowUser(sortedRecords)

	cost[0] = costs.NewCost(GeneralRatesStandard, GeneralRatesLow)

	return cost
}
