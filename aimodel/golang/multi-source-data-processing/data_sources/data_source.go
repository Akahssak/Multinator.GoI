package data_sources

// DataSource interface definition
type DataSource interface {
	FetchData() ([]map[string]interface{}, error)
}
