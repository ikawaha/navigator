package navigator

import (
	"time"
)

// A Log represents a single log message
type Log struct {
	// Message is the text of the log message.
	Message string
	// Location is the code location of the log message, if present
	Location string
	// Level is the log level ("DEBUG", "INFO", "WARNING", or "SEVERE").
	Level string
	// Time is the time the message was logged.
	Time time.Time
}
