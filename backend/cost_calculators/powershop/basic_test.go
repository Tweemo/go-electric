package cost_calculators

import (
	"testing"

	"github.com/tweemo/go-electric/rates"
	"github.com/tweemo/go-electric/utils"
)

// TestCalculateBasicCost prices a weekday peak record and a weekend off-peak record,
// confirming peak/off-peak split applies uniformly across day types and the daily
// charge covers every day.
func TestCalculateBasicCost(t *testing.T) {
	records, err := utils.CalculateDayPower([][]string{
		// Monday (weekday) 10:00 -> peak (7am-9pm), 2.0 kWh
		{"10/03/2025 10:00:00", "10/03/2025 10:30:00", "2.0"},
		// Saturday (weekend) 23:00 -> off-peak (9pm-7am), 3.0 kWh
		{"08/03/2025 23:00:00", "08/03/2025 23:30:00", "3.0"},
	})
	if err != nil {
		t.Fatal(err)
	}

	rate := rates.Rate{Pwh7amTo9pm: 0.30, Pwh9pmTo7am: 0.20, Daily: 2.00}

	// peak 2.0*0.30 + off-peak 3.0*0.20 + daily 2.00*2 days = 0.6 + 0.6 + 4.0 = 5.20
	got, err := CalculateBasicCost(records, rate, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 5.20 {
		t.Errorf("CalculateBasicCost = %v, want 5.20", got)
	}
}
