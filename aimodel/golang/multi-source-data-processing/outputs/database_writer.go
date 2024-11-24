package outputs

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func WriteToDatabase(data []map[string]interface{}, dbDSN string, tableName string) error {
	if len(data) == 0 {
		return nil
	}

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create placeholders and columns from first record
	var columns []string
	for column := range data[0] {
		columns = append(columns, column)
	}

	// Prepare statement
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert records
	for _, record := range data {
		values := make([]interface{}, len(columns))
		for i, col := range columns {
			values[i] = record[col]
		}

		if _, err := stmt.Exec(values...); err != nil {
			return err
		}
	}

	return nil
}
