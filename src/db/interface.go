package db

// LogRow represents a log entry
type LogRow struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Processed bool   `json:"processed"`
}

// DB interface defines the methods that any DB implementation should have
type DB interface {
	GetUnprocessedLogs(table string, markedColumn string, unmarkedValue string) ([]map[string]interface{}, error)
	MarkLogAsProcessed(table string, logID int, markedColumn string, markedValue string) error
}
