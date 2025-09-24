// Package log provides SOLID-compliant logging components.
package log

// Level defines the different levels for log messages.
type Level int

// Defines the available log levels.
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// levelToString converts a LogLevel to its string representation.
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
