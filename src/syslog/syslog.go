package syslog

import (
	"fmt"
	"net"
	"time"
)

// SendSyslogMessageTCP sends a formatted RFC5424 syslog message to the specified server over TCP.
func SendSyslogMessageTCP(syslogServer string, syslogPort int, syslogTransport string, sdata string, hostname string, appName string, procID string, msgID string, message string) (string, error) {
	// Create a TCP connection to the syslog server
	serverAddress := fmt.Sprintf("%s:%d", syslogServer, syslogPort)
	conn, err := net.Dial(syslogTransport, serverAddress)
	if err != nil {
		return "", fmt.Errorf("error connecting to syslog server: %w", err)
	}
	defer conn.Close()

	// Define syslog parameters
	facility := 4 // user-level messages
	severity := 2 // critical error
	priority := (facility * 8) + severity

	// Construct the PRI field (facility * 8 + severity)
	priField := fmt.Sprintf("<%d>1", priority)

	// Construct timestamp in ISO 8601 format
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	// Adding Structured Data (SDATA)
	sdataPart := fmt.Sprintf("[%s]", sdata)

	// Form the entire syslog message
	syslogMessage := fmt.Sprintf("%s %s %s %s %s %s %s %s\n", priField, timestamp, hostname, appName, procID, msgID, sdataPart, message)

	// Send the message via TCP to the syslog server
	_, err = conn.Write([]byte(syslogMessage))
	if err != nil {
		return "", fmt.Errorf("error sending syslog message: %w", err)
	}

	// Return the syslog message sent (for logging or debugging purposes)
	return syslogMessage, nil
}
