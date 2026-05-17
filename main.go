package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tweemo/go-electric/cost_calculators"
	"github.com/tweemo/go-electric/utils"
)

type RequestData struct {
	UsageData [][]string `json:"usage_data"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/costs", Costs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func Costs(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20)

	// Retrieve the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file in the data directory
	dataFilePath := filepath.Join("data/", "data.csv")
	dst, err := os.Create(dataFilePath)
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file's content to the new file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	// Get usage data and calculate power
	data := utils.GetUsageData(dataFilePath)
	sortedRecords := utils.CalculateDayPower(data)

	response := cost_calculators.AllPrices(sortedRecords)

	// Set content type to JSON and send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
