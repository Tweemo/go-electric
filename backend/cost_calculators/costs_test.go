package cost_calculators

import (
	"testing"

	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

func testRates() *rates.Rates {
	return &rates.Rates{Companies: []rates.Company{
		{Name: "Contact", Levy: 0.002, Plans: []rates.Plan{
			{Name: "GoodCharge", Standard: rates.Rate{Pwh7amTo9pm: 0.30, Pwh9pmTo7am: 0.15, Daily: 2.99}, Low: rates.Rate{Pwh7amTo9pm: 0.37, Pwh9pmTo7am: 0.18, Daily: 1.72}},
			{Name: "GoodNights", Standard: rates.Rate{Pwh: 0.32, Daily: 2.99}, Low: rates.Rate{Pwh: 0.39, Daily: 1.72}},
			{Name: "GoodWeekends", Standard: rates.Rate{Pwh: 0.29, Daily: 2.98}, Low: rates.Rate{Pwh: 0.35, Daily: 1.72}},
		}},
		{Name: "Nova", Levy: 0.003, Plans: []rates.Plan{
			{Name: "Basic", Standard: rates.Rate{Pwh: 0.25, Daily: 3.0}, Low: rates.Rate{Pwh: 0.31, Daily: 1.72}},
		}},
		{Name: "Powershop", Levy: 0, Plans: []rates.Plan{
			{Name: "Basic", Standard: rates.Rate{Pwh7amTo9pm: 0.33, Pwh9pmTo7am: 0.21, Daily: 2.75}, Low: rates.Rate{Pwh7amTo9pm: 0.37, Pwh9pmTo7am: 0.25, Daily: 1.95}},
		}},
	}}
}

// oneDayRecords builds a single day (09/03/2025) with 1.0 kWh of usage at hour 0.
func oneDayRecords(t *testing.T) []utils.DayPower {
	t.Helper()
	records, err := utils.CalculateDayPower([][]string{
		{"09/03/2025 00:00:00", "09/03/2025 00:30:00", "1.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return records
}

func TestAllPricesKnownValue(t *testing.T) {
	out, err := AllPrices(oneDayRecords(t), testRates())
	if err != nil {
		t.Fatal(err)
	}

	// Nova General Rates standard: (pwh + levy) * usage + daily * days
	// = (0.25 + 0.003) * 1.0 + 3.0 * 1 = 3.253 -> 3.25
	got := out["nova"]["GeneralRatesStandard"]
	if got != 3.25 {
		t.Errorf("GeneralRatesStandard = %v, want 3.25", got)
	}

	// Powershop Basic standard: the fixture's 1.0 kWh at hour 0 is off-peak.
	// off-peak 1.0 * 0.21 + daily 2.75 * 1 day = 2.96
	if got := out["powershop"]["BasicStandard"]; got != 2.96 {
		t.Errorf("BasicStandard = %v, want 2.96", got)
	}
}

func TestAllPricesHasAllKeys(t *testing.T) {
	out, err := AllPrices(oneDayRecords(t), testRates())
	if err != nil {
		t.Fatal(err)
	}

	wantContact := []string{
		"GoodChargeStandard", "GoodChargeLow", "GoodNightsStandard", "GoodNightsLow",
		"GoodWeekendsStandard", "GoodWeekendsLow", "SimpleRatesStandard", "SimpleRatesLow",
	}
	for _, k := range wantContact {
		if _, ok := out["contact"][k]; !ok {
			t.Errorf("missing contact key %q", k)
		}
	}
	for _, k := range []string{"GeneralRatesStandard", "GeneralRatesLow"} {
		if _, ok := out["nova"][k]; !ok {
			t.Errorf("missing nova key %q", k)
		}
	}
	for _, k := range []string{"BasicStandard", "BasicLow"} {
		if _, ok := out["powershop"][k]; !ok {
			t.Errorf("missing powershop key %q", k)
		}
	}
}

func TestAllPricesMissingPlanErrors(t *testing.T) {
	// Rates without the Contact plans the specs require should error, not panic.
	incomplete := &rates.Rates{Companies: []rates.Company{
		{Name: "Nova", Levy: 0.003, Plans: []rates.Plan{{Name: "Basic"}}},
	}}
	if _, err := AllPrices(oneDayRecords(t), incomplete); err == nil {
		t.Error("expected error when required plans are missing")
	}
}
