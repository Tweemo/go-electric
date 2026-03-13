package cost_calculators

type Cost struct {
	standard, low float64
}

func NewCost(standard float64, low float64) Cost {
	return Cost{
		standard: standard,
		low:      low,
	}
}
