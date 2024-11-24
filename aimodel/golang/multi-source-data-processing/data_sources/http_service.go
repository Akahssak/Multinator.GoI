package data_sources

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPServiceSource struct {
	URL string
}

// FetchData fetches data from an HTTP service
func (h *HTTPServiceSource) FetchData() ([]map[string]interface{}, error) {
	resp, err := http.Get(h.URL)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Decode JSON response
	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return data, nil
}
