package main

import (
	"db2syslog/db"
	"db2syslog/syslog"
	"db2syslog/utils"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql" // MySQL Driver
	_ "github.com/sijms/go-ora/v2"     // go-ora Driver for Oracle
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// Define the command-line flag for the config file
	configFile := pflag.String("config.file", "config.yml", "Path to the configuration file")
	pflag.Parse()

	// Load configuration from the config file
	utils.InitConfig(*configFile)

	// Get the database type from the config file (mysql or oracle)
	dbType := viper.GetString("database.type")
	table := viper.GetString("database.table")
	markedColumn := viper.GetString("database.markedcolumn")
	markedValue := viper.GetString("database.markedvalue")
	unmarkedValue := viper.GetString("database.unmarkedvalue")
	idColumn := viper.GetString("database.keycolumn")

	// Set up the database connection based on the config
	var database db.DB
	var err error
	switch dbType {
	case "mysql":
		database, err = db.NewMySQLConnection(
			viper.GetString("database.host"),
			viper.GetInt("database.port"),
			viper.GetString("database.username"),
			viper.GetString("database.password"),
			viper.GetString("database.dbname"),
		)
	case "oracle":
		database, err = db.NewOracleConnection(
			viper.GetString("database.host"),
			viper.GetInt("database.port"),
			viper.GetString("database.username"),
			viper.GetString("database.password"),
			viper.GetString("database.dbname"),
			viper.GetStringMapString("database.presets"),
		)
	default:
		log.Fatal("Unsupported DB type in config file")
	}
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
		os.Exit(1)
	}
	// Ensure we close the DB connection properly
	defer func() {
		switch v := database.(type) {
		case *db.MySQLDB:
			v.DB.Close() // Close MySQL connection
		case *db.OracleDB:
			v.DB.Close() // Close Oracle connection
		default:
			log.Fatalf("Unknown database type: %T", v)
		}
	}()

	// Example: Get unprocessed logs and send them to syslog
	logRows, err := database.GetUnprocessedLogs(table, markedColumn, unmarkedValue)
	if err != nil {
		log.Fatal("Error fetching unprocessed logs: ", err)
		os.Exit(1)
	}

	// Parse dropcolumn as a list (slice)
	dropColumnList := viper.GetStringSlice("syslog.dropcolumn")

	for _, logRow := range logRows {
		// Debugging: Print the logRow content for troubleshooting
		log.Printf("Processing logRow: %v", logRow)

		// Format logRow map as a JSON-like string for Syslog
		logData := "{"
		for column, value := range logRow {
			column = strings.ToLower(column)
			//drops Columns from Json
			if utils.StringInSlice(column, dropColumnList) {
				continue
			}
			if column != markedColumn { // Exclude 'processed' field
				logData += fmt.Sprintf(`"%s": "%s",`, column, utils.ConvertToString(value))
			}
		}
		// Remove the last comma
		if len(logData) > 1 {
			logData = logData[:len(logData)-1]
		}
		logData += "}"

		// Send log to Syslog
		syslogMessage, err := syslog.SendSyslogMessageTCP(
			viper.GetString("syslog.server"),
			viper.GetInt("syslog.port"),
			viper.GetString("syslog.transport"),
			viper.GetString("syslog.sdata"),
			viper.GetString("syslog.hostname"),
			viper.GetString("syslog.app_name"),
			viper.GetString("syslog.proc_id"),
			viper.GetString("syslog.msg_id"),
			logData,
		)
		if err != nil {
			// Skip marking the log as processed on syslog error
			log.Printf("Error sending syslog message: %v", err)
			continue
		}
		log.Printf("Sent Syslog Message: %s\n", syslogMessage)

		// Ensure that we are accessing the correct key to mark as "processed"
		id, err := utils.ConvertToInt(logRow[idColumn])
		if err != nil {
			log.Printf("Error: id field conversion failed: %v", err)
			continue
		}

		// Mark the log as processed in the database
		err = database.MarkLogAsProcessed(table, id, markedColumn, markedValue)
		if err != nil {
			log.Printf("Error updating processed flag for log %s %d: %v", idColumn, id, err)
			os.Exit(1)
		}
	}

}
