package transformations

func EnhanceData(data []map[string]interface{}, enhancementFunc func(map[string]interface{}) map[string]interface{}) []map[string]interface{} {
	enhancedData := make([]map[string]interface{}, len(data))
	for i, record := range data {
		enhancedData[i] = enhancementFunc(record)
	}
	return enhancedData
}
