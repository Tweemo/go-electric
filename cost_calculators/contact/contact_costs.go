package cost_calculators

import (
	costs "github.com/tweemo/go-electric/cost_calculators"
	"github.com/tweemo/go-electric/utils"
)

func ContactCosts(sortedRecords []utils.DayPower) [4]costs.Cost {
	var cost [4]costs.Cost

	// GoodCharge
	GoodChargeStandard := GoodChargeStandardUser(sortedRecords)
	GoodChargeLow := GoodChargeLowUser(sortedRecords)

	cost[0] = costs.NewCost(GoodChargeStandard, GoodChargeLow)

	// GoodNights
	GoodNightsStandard := GoodNightsStandardUser(sortedRecords)
	GoodNightsLow := GoodNightsLowUser(sortedRecords)

	cost[1] = costs.NewCost(GoodNightsStandard, GoodNightsLow)

	// GoodWeekends
	GoodWeekendsStandard := GoodWeekendsStandardUser(sortedRecords)
	GoodWeekendsLow := GoodWeekendsLowUser(sortedRecords)

	cost[2] = costs.NewCost(GoodWeekendsStandard, GoodWeekendsLow)

	// SimpleRates
	SimpleRatesStandard := SimpleRatesStandardUser(sortedRecords)
	SimpleRatesLow := SimpleRatesLowUser(sortedRecords)

	cost[3] = costs.NewCost(SimpleRatesStandard, SimpleRatesLow)

	return cost
}
