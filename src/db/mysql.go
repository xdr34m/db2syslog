package db

import (
	"database/sql"
	"fmt"

	// MySQL Driver
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB struct implements the DB interface for MySQL
type MySQLDB struct {
	DB *sql.DB
}

// NewMySQLConnection creates and returns a new MySQLDB instance
func NewMySQLConnection(host string, port int, username string, password string, dbname string) (*MySQLDB, error) {
	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Return a new MySQLDB instance
	return &MySQLDB{DB: db}, nil
}

// GetUnprocessedLogs retrieves unprocessed logs from MySQL
func (m *MySQLDB) GetUnprocessedLogs(table string, markedColumn string, unmarkedValue string) ([]map[string]interface{}, error) {
	// Construct the query to fetch all columns (except processed)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", table, markedColumn, unmarkedValue)

	// Execute the query
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Get column names dynamically
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch column names: %w", err)
	}

	var logs []map[string]interface{}
	for rows.Next() {
		// Create a slice of pointers to hold the row values
		values := make([]interface{}, len(columns))
		// Create a slice to scan each column into
		columnPointers := make([]interface{}, len(columns))
		for i := range columnPointers {
			columnPointers[i] = &values[i]
		}

		// Scan the row into the pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Create a map to store the column name and value
		rowMap := make(map[string]interface{})
		for i, column := range columns {
			// Skip the 'processed' column
			if column == "processed" {
				continue
			}
			rowMap[column] = values[i]
		}

		logs = append(logs, rowMap)
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return logs, nil
}

// MarkLogAsProcessed updates the processed flag in MySQL
func (m *MySQLDB) MarkLogAsProcessed(table string, logID int, markedColumn string, markedValue string) error {
	query := fmt.Sprintf("UPDATE %s SET %s = %s WHERE id = ?", table, markedColumn, markedValue)
	_, err := m.DB.Exec(query, logID)
	if err != nil {
		return err
	}
	return nil
}
