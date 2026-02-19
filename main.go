package main

import (
	"fmt"

	contact_cost_calculators "github.com/tweemo/go-electric/cost_calculators/contact"
	"github.com/tweemo/go-electric/utils"
)

func main() {
	usageData := utils.GetUsageData()
	sortedRecords := utils.CalculateDayPower(usageData)
	contact_good_charge_su := contact_cost_calculators.CalculateGoodChargeStandardUserCost(sortedRecords)
	fmt.Println(contact_good_charge_su)
}
