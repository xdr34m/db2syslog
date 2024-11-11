package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// go-ora Driver for Oracle
	_ "github.com/sijms/go-ora/v2"
)

// OracleDB struct implements the DB interface for Oracle
type OracleDB struct {
	DB *sql.DB
}

// NewOracleConnection creates and returns a new OracleDB instance
func NewOracleConnection(host string, port int, username string, password string, serviceName string) (*OracleDB, error) {
	// Construct the DSN (Data Source Name) for go-ora
	dsn := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", username, password, host, port, serviceName)
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	// Test the database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	// Set the session's time zone to 'Europe/Berlin' (German time zone)
	_, err = db.Exec("ALTER SESSION SET TIME_ZONE = 'Europe/Berlin'")
	if err != nil {
		return nil, fmt.Errorf("error setting session time zone: %w", err)
	}

	// Set NLS parameters to match German conventions for date/time formats
	_, err = db.Exec("ALTER SESSION SET NLS_DATE_FORMAT = 'DD.MM.YYYY HH24:MI:SS'")
	if err != nil {
		return nil, fmt.Errorf("error setting NLS_DATE_FORMAT: %w", err)
	}
	_, err = db.Exec("ALTER SESSION SET NLS_TIMESTAMP_FORMAT = 'DD.MM.YYYY HH24:MI:SS.FF'")
	if err != nil {
		return nil, fmt.Errorf("error setting NLS_TIMESTAMP_FORMAT: %w", err)
	}

	// Set the NLS language and territory to German
	_, err = db.Exec("ALTER SESSION SET NLS_LANGUAGE = 'GERMAN'")
	if err != nil {
		return nil, fmt.Errorf("error setting NLS_LANGUAGE: %w", err)
	}
	_, err = db.Exec("ALTER SESSION SET NLS_TERRITORY = 'GERMANY'")
	if err != nil {
		return nil, fmt.Errorf("error setting NLS_TERRITORY: %w", err)
	}

	// Return a new OracleDB instance
	log.Println("Connected to Oracle database successfully.")
	return &OracleDB{DB: db}, nil
}

// GetUnprocessedLogs retrieves unprocessed logs from the Oracle table
func (o *OracleDB) GetUnprocessedLogs(table string, markedColumn string, unmarkedValue string) ([]map[string]interface{}, error) {
	// Construct the query to fetch all columns except processed
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", table, markedColumn, unmarkedValue)

	// Execute the query
	rows, err := o.DB.Query(query)
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
		// Create a slice to hold the row values
		values := make([]interface{}, len(columns))
		// Create pointers to each value
		columnPointers := make([]interface{}, len(columns))
		for i := range columnPointers {
			columnPointers[i] = &values[i]
		}

		// Scan the row into the pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Create a map to store column name and value pairs
		rowMap := make(map[string]interface{})
		for i, column := range columns {
			column = strings.ToLower(column)
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

// MarkLogAsProcessed updates the processed flag in the Oracle table
func (o *OracleDB) MarkLogAsProcessed(table string, logID int, markedColumn string, markedValue string) error {
	query := fmt.Sprintf("UPDATE %s SET %s = %s WHERE id = :1", table, markedColumn, markedValue)
	_, err := o.DB.Exec(query, logID)
	if err != nil {
		return fmt.Errorf("error updating log as processed: %w", err)
	}
	return nil
}
