package model

import "time"

// LogEntry represents one log line after parsing.
type LogEntry struct {
	ID uint64

	Timestamp time.Time

	Namespace string
	Pod       string

	Raw string

	Level         string
	Message       string
	TraceID       string
	CorrelationID string
	RequestID     string

	// All extracted string fields from the JSON.
	Fields map[string]string
}
