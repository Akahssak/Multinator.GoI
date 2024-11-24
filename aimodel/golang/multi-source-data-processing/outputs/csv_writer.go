package outputs

import (
	"encoding/csv"
	"fmt"
	"os"
)

func WriteToCSV(data []map[string]interface{}, outputPath string) error {
	// Only proceed if data is not empty
	if len(data) == 0 {
		return fmt.Errorf("no data to write to CSV")
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Get headers from the first record
	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}

	// Write the header row
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("could not write header to CSV: %v", err)
	}

	// Write records
	for _, record := range data {
		row := make([]string, len(headers))
		for i, header := range headers {
			if value, exists := record[header]; exists {
				row[i] = fmt.Sprintf("%v", value)
			}
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("could not write row to CSV: %v", err)
		}
	}

	return nil
}
