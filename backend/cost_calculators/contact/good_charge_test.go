package cost_calculators

import (
	"testing"

	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// TestCalculateGoodChargeCostWeekendMorning is a regression test for the weekend
// off-peak morning (00:00-07:00) window. It places usage in the 0-7 range on both a
// weekday and a weekend day; the old code mistakenly priced weekday usage twice and
// dropped weekend morning usage entirely.
func TestCalculateGoodChargeCostWeekendMorning(t *testing.T) {
	records, err := utils.CalculateDayPower([][]string{
		// Saturday (weekend) 03:00 -> off-peak morning, 4.0 kWh
		{"08/03/2025 03:00:00", "08/03/2025 03:30:00", "4.0"},
		// Monday (weekday) 03:00 -> off-peak morning, 1.0 kWh
		{"10/03/2025 03:00:00", "10/03/2025 03:30:00", "1.0"},
	})
	if err != nil {
		t.Fatal(err)
	}

	rate := rates.Rate{Pwh7amTo9pm: 0.30, Pwh9pmTo7am: 0.15, Daily: 2.00}

	// All 5.0 kWh is off-peak: (4.0 weekend + 1.0 weekday) * 0.15 + daily 2.00*2 days
	// = 0.75 + 4.00 = 4.75. The old buggy line would yield 4.30.
	got, err := CalculateGoodChargeCost(records, rate, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 4.75 {
		t.Errorf("CalculateGoodChargeCost = %v, want 4.75", got)
	}
}
