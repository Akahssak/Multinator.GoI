package transformations

import "errors"

func AggregateData(data []map[string]interface{}, field string, operation string) (float64, error) {
	var sum float64
	count := 0

	for _, record := range data {
		value, exists := record[field]
		if !exists {
			continue
		}

		numValue, ok := value.(float64)
		if !ok {
			return 0, errors.New("field value is not a number")
		}

		sum += numValue
		count++
	}

	switch operation {
	case "sum":
		return sum, nil
	case "average":
		if count == 0 {
			return 0, errors.New("no valid records to aggregate")
		}
		return sum / float64(count), nil
	default:
		return 0, errors.New("unsupported operation")
	}
}
