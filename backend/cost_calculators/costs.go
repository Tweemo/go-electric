package cost_calculators

import (
	contact "github.com/tweemo/go-electric/cost_calculators/contact"
	nova "github.com/tweemo/go-electric/cost_calculators/nova"
	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// calcFn computes a plan's total cost from usage records, its resolved rate, and
// the company levy (zero when the plan does not apply a levy).
type calcFn func(records []utils.DayPower, rate rates.Rate, levy float64) (float64, error)

// planSpec describes one priced entry in the response.
type planSpec struct {
	group   string // top-level response key, e.g. "contact"
	key     string // plan key, e.g. "GoodChargeStandard"
	company string // company name in rates.json
	plan    string // plan name in rates.json
	tier    string // "standard" or "low"
	levyOf  string // company to source the levy from; "" means no levy
	calc    calcFn
}

// NOTE: SimpleRates intentionally resolves the "GoodNights" plan to preserve the
// existing behavior — this looks like a bug worth revisiting separately.
var specs = []planSpec{
	{"contact", "GoodChargeStandard", "Contact", "GoodCharge", "standard", "", contact.CalculateGoodChargeCost},
	{"contact", "GoodChargeLow", "Contact", "GoodCharge", "low", "", contact.CalculateGoodChargeCost},
	{"contact", "GoodNightsStandard", "Contact", "GoodNights", "standard", "", contact.CalculateGoodNightsCost},
	{"contact", "GoodNightsLow", "Contact", "GoodNights", "low", "", contact.CalculateGoodNightsCost},
	{"contact", "GoodWeekendsStandard", "Contact", "GoodWeekends", "standard", "", contact.CalculateGoodWeekendsCost},
	{"contact", "GoodWeekendsLow", "Contact", "GoodWeekends", "low", "", contact.CalculateGoodWeekendsCost},
	{"contact", "SimpleRatesStandard", "Contact", "GoodNights", "standard", "Contact", contact.CalculateSimpleRatesCost},
	{"contact", "SimpleRatesLow", "Contact", "GoodNights", "low", "Contact", contact.CalculateSimpleRatesCost},
	{"nova", "GeneralRatesStandard", "Nova", "Basic", "standard", "Nova", nova.CalculateGeneralRatesCost},
	{"nova", "GeneralRatesLow", "Nova", "Basic", "low", "Nova", nova.CalculateGeneralRatesCost},
}

// AllPrices computes every plan's cost for the given usage records using the
// provided rates table. It returns an error rather than crashing on bad data.
func AllPrices(records []utils.DayPower, r *rates.Rates) (map[string]map[string]float64, error) {
	out := map[string]map[string]float64{}

	for _, s := range specs {
		rate, err := r.Get(s.company, s.plan, s.tier)
		if err != nil {
			return nil, err
		}

		var levy float64
		if s.levyOf != "" {
			if levy, err = r.Levy(s.levyOf); err != nil {
				return nil, err
			}
		}

		cost, err := s.calc(records, rate, levy)
		if err != nil {
			return nil, err
		}

		if out[s.group] == nil {
			out[s.group] = map[string]float64{}
		}
		out[s.group][s.key] = cost
	}

	return out, nil
}
