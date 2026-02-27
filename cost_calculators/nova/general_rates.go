package cost_calculatr

import "github.com/tweemo/go-electric/utils"

// Todo this could realistically just be a single func that we pass low or standard to
func NovaGeneralRatesStandardUser(sortedRecords []utils.DayPower) float64 {
	standard := utils.GetRate("Nova", "Basic", "standard")
	standardRateMap := standard.(map[string]interface{})
	pwh, _ := utils.GetFloat(standardRateMap["pwh"])
	daily, _ := utils.GetFloat(standardRateMap["daily"])

	totalCost := CalculateGeneralRatesCost(sortedRecords, pwh, daily)
	return totalCost
}

func NovaGeneralRatesLowUser(sortedRecords []utils.DayPower) float64 {
	low := utils.GetRate("Nova", "Basic", "low")
	lowRateMap := low.(map[string]interface{})
	pwh, _ := utils.GetFloat(lowRateMap["pwh"])
	daily, _ := utils.GetFloat(lowRateMap["daily"])

	totalCost := CalculateGeneralRatesCost(sortedRecords, pwh, daily)
	return totalCost
}

func CalculateGeneralRatesCost(sortedRecords []utils.DayPower, pwh float64, daily float64) float64 {
	usage := utils.TotalUsage(sortedRecords)
	levy := utils.GetLevy("Nova")

	cost := (pwh + levy) * usage
	cost += daily * 30

	roundedCost, err := utils.RoundFloat(cost, 2)
	if err != nil {
		panic(err)
	}

	return roundedCost
}
