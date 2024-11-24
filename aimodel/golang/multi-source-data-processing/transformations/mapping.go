package transformations

func MapFields(data []map[string]interface{}, mappingRules map[string]string) []map[string]interface{} {
	var mappedData []map[string]interface{}
	for _, record := range data {
		mappedRecord := make(map[string]interface{})
		// Copy unmapped fields
		for k, v := range record {
			mappedRecord[k] = v
		}
		// Apply mapping rules
		for oldKey, newKey := range mappingRules {
			if value, exists := record[oldKey]; exists {
				mappedRecord[newKey] = value
				delete(mappedRecord, oldKey)
			}
		}
		mappedData = append(mappedData, mappedRecord)
	}
	return mappedData
}
