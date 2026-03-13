package main

import (
	"fmt"

	contact_cost_calculator "github.com/tweemo/go-electric/cost_calculators/contact"
	nova_cost_calculator "github.com/tweemo/go-electric/cost_calculators/nova"
	"github.com/tweemo/go-electric/utils"
)

func main() {
	usageData := utils.GetUsageData()
	sortedRecords := utils.CalculateDayPower(usageData)

	ContactCost := contact_cost_calculator.ContactCosts(sortedRecords)
	NovaCost := nova_cost_calculator.NovaCosts(sortedRecords)
	fmt.Println(ContactCost, NovaCost)
}
