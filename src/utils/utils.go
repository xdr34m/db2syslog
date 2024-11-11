package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/viper"
)

// Function to check if a string is in a slice
func StringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
func InitConfig(configFile string) {
	// Set the file name and path for the config file
	viper.SetConfigFile(configFile)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

// ConvertToString converts various types of values to strings.
func ConvertToString(value interface{}) string {
	switch v := value.(type) {
	case []byte:
		return string(v) // If it's a byte slice, convert to string
	case int, int64:
		return fmt.Sprintf("%d", v) // For int or int64, convert to string
	case string:
		return v // If it's already a string, return it
	default:
		return fmt.Sprintf("%v", v) // For any other type, use default string representation
	}
}

// Convert various possible 'id' types to an integer
func ConvertToInt(value interface{}) (int, error) {
	if value == nil {
		return 0, fmt.Errorf("id value is nil")
	}

	switch v := value.(type) {
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		idInt64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error converting id from string to int64: %w", err)
		}
		return int(idInt64), nil
	default:
		return 0, fmt.Errorf("unsupported type for id field %T", v)
	}
}
