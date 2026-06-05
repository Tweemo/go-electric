// Package rates loads the power-company pricing table (data/rates.json) once into
// typed structs and answers rate/levy lookups. Loading happens at startup; lookups
// are pure in-memory reads that return errors instead of crashing the process.
package rates

import (
	"encoding/json"
	"fmt"
	"os"
)

// Rate is the pricing for one plan at one usage tier (standard/low).
type Rate struct {
	Pwh         float64 `json:"pwh"`
	Pwh7amTo9pm float64 `json:"pwh_7am_9pm"`
	Pwh9pmTo7am float64 `json:"pwh_9pm_7am"`
	Daily       float64 `json:"daily"`
}

// Plan is a named pricing plan with standard and low usage tiers.
type Plan struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Standard    Rate   `json:"standard"`
	Low         Rate   `json:"low"`
}

// Company is a power company with a levy and a set of plans.
type Company struct {
	Name  string  `json:"name"`
	Levy  float64 `json:"levy"`
	Plans []Plan  `json:"plans"`
}

// Rates is the full pricing table.
type Rates struct {
	Companies []Company `json:"power_companies"`
}

// Load reads and parses the rates table from path. Call once at startup.
func Load(path string) (*Rates, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rates file %q: %w", path, err)
	}

	var r Rates
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("parse rates file %q: %w", path, err)
	}
	if len(r.Companies) == 0 {
		return nil, fmt.Errorf("rates file %q contains no companies", path)
	}
	return &r, nil
}

func (r *Rates) company(name string) (Company, error) {
	for _, c := range r.Companies {
		if c.Name == name {
			return c, nil
		}
	}
	return Company{}, fmt.Errorf("no rates for company %q", name)
}

// Get returns the rate for a company/plan at the given usage tier
// ("standard" or "low").
func (r *Rates) Get(company, plan, usageType string) (Rate, error) {
	c, err := r.company(company)
	if err != nil {
		return Rate{}, err
	}
	for _, p := range c.Plans {
		if p.Name != plan {
			continue
		}
		switch usageType {
		case "standard":
			return p.Standard, nil
		case "low":
			return p.Low, nil
		default:
			return Rate{}, fmt.Errorf("unknown usage type %q", usageType)
		}
	}
	return Rate{}, fmt.Errorf("no plan %q for company %q", plan, company)
}

// Levy returns the regulatory levy for a company.
func (r *Rates) Levy(company string) (float64, error) {
	c, err := r.company(company)
	if err != nil {
		return 0, err
	}
	return c.Levy, nil
}
