package cost_calculators

import (
	contact_cost_calculator "github.com/tweemo/go-electric/cost_calculators/contact"
	nova_cost_calculator "github.com/tweemo/go-electric/cost_calculators/nova"
	"github.com/tweemo/go-electric/utils"
)

func AllPrices(sortedRecords []utils.DayPower) map[string]map[string]float64 {
	return map[string]map[string]float64{
		"contact": {
			"GoodChargeStandard":   contact_cost_calculator.GoodChargeStandardUser(sortedRecords),
			"GoodChargeLow":        contact_cost_calculator.GoodChargeLowUser(sortedRecords),
			"GoodNightsStandard":   contact_cost_calculator.GoodNightsStandardUser(sortedRecords),
			"GoodNightsLow":        contact_cost_calculator.GoodNightsLowUser(sortedRecords),
			"GoodWeekendsStandard": contact_cost_calculator.GoodWeekendsStandardUser(sortedRecords),
			"GoodWeekendsLow":      contact_cost_calculator.GoodWeekendsLowUser(sortedRecords),
			"SimpleRatesStandard":  contact_cost_calculator.SimpleRatesStandardUser(sortedRecords),
			"SimpleRatesLow":       contact_cost_calculator.SimpleRatesLowUser(sortedRecords),
		},
		"nova": {
			"GeneralRatesStandard": nova_cost_calculator.NovaGeneralRatesStandardUser(sortedRecords),
			"GeneralRatesLow":      nova_cost_calculator.NovaGeneralRatesLowUser(sortedRecords),
		},
	}
}
