package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

// MustFloat64Env returns the float64 value of an environment variable, or panics if it's not set
func MustFloat64Env(key string) float64 {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("env %s is required and must be set", key)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("env %s must be a valid number, got %q: %v", key, s, err)
	}
	return v
}

func GetRate(company string, plan string, usageType string) interface{} {
	// Load the JSON file
	file, err := os.Open("data/rates.json")
	if err != nil {
		log.Fatalf("Failed to open rates file: %v", err)
	}
	defer file.Close()

	// Decode the JSON data
	var ratesData map[string]interface{}

	if err := json.NewDecoder(file).Decode(&ratesData); err != nil {
		log.Fatalf("Failed to decode rates file: %v", err)
	}

	// Find the company in the rates data
	var companyData map[string]interface{}
	for _, companyEntry := range ratesData["power_companies"].([]interface{}) {
		companyEntryMap := companyEntry.(map[string]interface{})
		if companyEntryMap["name"] == company {
			companyData = companyEntryMap
			break
		}
	}

	if companyData == nil {
		log.Fatalf("No rates found for company %s", company)
	}

	// Find the plan in the company data
	var planData map[string]interface{}
	for _, planEntry := range companyData["plans"].([]interface{}) {
		planEntryMap := planEntry.(map[string]interface{})
		if planEntryMap["name"] == plan {
			planData = planEntryMap
			break
		}
	}

	if planData == nil {
		log.Fatalf("No rates found for plan %s in company %s", plan, company)
	}

	// Todo fix this so we send back a float rather than interface and all usages dont need to run GetFloat
	return planData[usageType]
}

func GetFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	if !v.Type().ConvertibleTo(reflect.TypeOf(float64(0))) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	return v.Convert(reflect.TypeOf(float64(0))).Float(), nil
}

func GetLevy(company string) float64 {
	// Load the JSON file
	file, err := os.Open("data/rates.json")
	if err != nil {
		log.Fatalf("Failed to open rates file: %v", err)
	}
	defer file.Close()

	// Decode the JSON data
	var ratesData map[string]interface{}

	if err := json.NewDecoder(file).Decode(&ratesData); err != nil {
		log.Fatalf("Failed to decode rates file: %v", err)
	}

	// Find the company in the rates data
	var companyData map[string]interface{}
	for _, companyEntry := range ratesData["power_companies"].([]interface{}) {
		companyEntryMap := companyEntry.(map[string]interface{})
		if companyEntryMap["name"] == company {
			companyData = companyEntryMap
			break
		}
	}

	if companyData == nil {
		log.Fatalf("No rates found for company %s", company)
	}

	levy, _ := GetFloat(companyData["levy"])

	return levy
}
