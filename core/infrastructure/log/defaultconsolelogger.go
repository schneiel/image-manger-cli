// Package log provides SOLID-compliant logging components.
package log

import (
	"errors"
	"fmt"
	"log"
)

// DefaultConsoleLogger logs messages to the standard output.
type DefaultConsoleLogger struct {
	stdoutLogger  *log.Logger
	stderrLogger  *log.Logger
	minLevel      Level
	consoleWriter ConsoleWriter
}

// NewDefaultConsoleLogger creates a new logger that writes to the console.
func NewDefaultConsoleLogger(minLevel Level) (*DefaultConsoleLogger, error) {
	return NewDefaultConsoleLoggerWithWriter(minLevel, NewDefaultConsoleWriter())
}

// NewDefaultConsoleLoggerWithWriter creates a new logger with injected console writer.
func NewDefaultConsoleLoggerWithWriter(minLevel Level, consoleWriter ConsoleWriter) (*DefaultConsoleLogger, error) {
	if consoleWriter == nil {
		return nil, errors.New("consoleWriter cannot be nil")
	}
	return &DefaultConsoleLogger{
		stdoutLogger:  log.New(consoleWriter.Stdout(), "", log.Ldate|log.Ltime),
		stderrLogger:  log.New(consoleWriter.Stderr(), "", log.Ldate|log.Ltime),
		minLevel:      minLevel,
		consoleWriter: consoleWriter,
	}, nil
}

// SetLevel configures the minimum log level for console output. Messages below this level are filtered out.
// Use DEBUG for development, INFO for production, WARN for important issues, ERROR for critical problems.
func (l *DefaultConsoleLogger) SetLevel(level Level) {
	l.minLevel = level
}

func (l *DefaultConsoleLogger) log(level Level, format string, v ...interface{}) {
	if level < l.minLevel {
		return
	}
	message := fmt.Sprintf(format, v...)

	// Use stderr for ERROR level, stdout for others
	if level == ERROR {
		l.stderrLogger.Printf("[%s] %s", level.String(), message)
	} else {
		l.stdoutLogger.Printf("[%s] %s", level.String(), message)
	}
}

// Public log methods

// Debug logs a debug level message.
func (l *DefaultConsoleLogger) Debug(message string) { l.log(DEBUG, "%s", message) }

// Info logs an info level message.
func (l *DefaultConsoleLogger) Info(message string) { l.log(INFO, "%s", message) }

// Warn logs a warning level message.
func (l *DefaultConsoleLogger) Warn(message string) { l.log(WARN, "%s", message) }

// Error logs an error level message.
func (l *DefaultConsoleLogger) Error(message string) { l.log(ERROR, "%s", message) }

// Public formatted log methods

// Debugf logs a formatted debug level message.
func (l *DefaultConsoleLogger) Debugf(format string, v ...interface{}) { l.log(DEBUG, format, v...) }

// Infof logs a formatted info level message.
func (l *DefaultConsoleLogger) Infof(format string, v ...interface{}) { l.log(INFO, format, v...) }

// Warnf logs a formatted warning level message.
func (l *DefaultConsoleLogger) Warnf(format string, v ...interface{}) { l.log(WARN, format, v...) }

// Errorf logs a formatted error level message.
func (l *DefaultConsoleLogger) Errorf(format string, v ...interface{}) { l.log(ERROR, format, v...) }
