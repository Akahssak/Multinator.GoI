package transformations

// IsAdult is an exported function for reuse
func IsAdult(record map[string]interface{}) bool {
	if age, ok := record["age"].(float64); ok {
		return age >= 18
	}
	return false
}

func FilterData(data []map[string]interface{}, condition func(map[string]interface{}) bool) []map[string]interface{} {
	var filteredData []map[string]interface{}
	for _, record := range data {
		if condition(record) {
			filteredData = append(filteredData, record)
		}
	}
	return filteredData
}
