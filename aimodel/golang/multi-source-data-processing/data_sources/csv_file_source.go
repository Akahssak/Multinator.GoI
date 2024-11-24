package data_sources

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xuri/excelize/v2" // Import for XLSX handling
)

type FileSource struct {
	FilePath string
}

// FetchData reads data from either a CSV or XLSX file based on the file extension
func (f *FileSource) FetchData() ([]map[string]interface{}, error) {
	// Determine the file type based on the extension
	if strings.HasSuffix(strings.ToLower(f.FilePath), ".csv") {
		return f.fetchCSVData()
	} else if strings.HasSuffix(strings.ToLower(f.FilePath), ".xlsx") {
		return f.fetchXLSXData()
	} else {
		return nil, fmt.Errorf("unsupported file type for file: %s", f.FilePath)
	}
}

// fetchCSVData reads data from a CSV file
func (f *FileSource) fetchCSVData() ([]map[string]interface{}, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, value := range record {
			row[headers[i]] = value
		}
		result = append(result, row)
	}

	return result, nil
}

// fetchXLSXData reads data from an XLSX file
func (f *FileSource) fetchXLSXData() ([]map[string]interface{}, error) {
	xlFile, err := excelize.OpenFile(f.FilePath)
	if err != nil {
		return nil, err
	}
	defer xlFile.Close()

	// Get the first sheet name
	sheetName := xlFile.GetSheetName(1)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in file: %s", f.FilePath)
	}

	// Read all rows in the sheet
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no data found in sheet: %s", sheetName)
	}

	// First row is the header
	headers := rows[0]

	var result []map[string]interface{}
	for _, row := range rows[1:] { // Skip header row
		rowData := make(map[string]interface{})
		for i, cell := range row {
			if i < len(headers) {
				rowData[headers[i]] = cell
			}
		}
		result = append(result, rowData)
	}

	return result, nil
}
