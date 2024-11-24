package main

import (
	"fmt"
	"log"
	"multi-source-data-processing/data_sources"
	"multi-source-data-processing/outputs"
)

func main() {
	// Initialize data sources
	sources := []data_sources.DataSource{
		&data_sources.FileSource{FilePath: "C:/Users/karthikeya.k/OneDrive/Documents/Desktop/golang/multi-source-data-processing/bank_marketing.csv"},
		&data_sources.JSONFileSource{FilePath: "C:/Users/karthikeya.k/OneDrive/Documents/Desktop/golang/multi-source-data-processing/json.json"},
		&data_sources.HTTPServiceSource{URL: "http://example.com/api/data"},
	}

	// Fetch and combine data from all sources
	var data []map[string]interface{}
	for _, source := range sources {
		sourceData, err := source.FetchData()
		if err != nil {
			log.Fatalf("Error fetching data from source: %v", err)
		}
		data = append(data, sourceData...)
	}

	// Only print fetched data if it's not empty
	if len(data) == 0 {
		log.Fatal("No data fetched from source.")
	}

	// Print fetched data to check the contents
	fmt.Printf("Fetched data: %+v\n", data)

	// Apply data transformations
	//mappingRules := map[string]string{"old_name": "new_name"} // Update field names as needed
	//data = transformations.MapFields(data, mappingRules)
	//data = transformations.FilterData(data, transformations.IsAdult)

	// Check if data is empty after transformations
	if len(data) == 0 {
		log.Fatal("Transformed data is empty.")
	}

	// Print transformed data only if it's not empty
	fmt.Printf("Transformed data: %+v\n", data)

	// Write output to CSV file
	outputFile := "output.csv"
	if err := outputs.WriteToCSV(data, outputFile); err != nil {
		log.Fatalf("Error writing data to CSV: %v", err)
	}

	fmt.Printf("Data processed and saved successfully to %s.\n", outputFile)
}
