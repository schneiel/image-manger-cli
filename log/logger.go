// Package log provides a simple leveled logger that can write to both the console and a file.
package log

import (
	"ImageManager/i18n"
	"fmt"
	"log"
	"os"
	"sync"
)

// LogLevel defines the different levels for log messages.
type LogLevel int

// Defines the available log levels.
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	// currentLogLevel sets the threshold for console output.
	currentLogLevel LogLevel = INFO
	fileLogger      *log.Logger
	logFile         *os.File
	logFilePath     = "application.log"
	// fileLoggerActive indicates whether the file logger is initialized.
	fileLoggerActive = false
	mu               sync.Mutex
)

// InitFileLogger opens the log file and initializes the fileLogger.
func InitFileLogger() error {
	mu.Lock()
	defer mu.Unlock()

	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Return the error to be handled in main(), where i18n is safely initialized.
		return fmt.Errorf(i18n.T("ErrorOpeningLogFile", map[string]interface{}{"Path": logFilePath, "Error": err}))
	}

	fileLogger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	fileLoggerActive = true
	return nil
}

// CloseFileLogger closes the log file if it is open.
func CloseFileLogger() {
	mu.Lock()
	defer mu.Unlock()
	if logFile != nil {
		logFile.Close()
	}
}

// SetLogLevel sets the log level for console output.
func SetLogLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	currentLogLevel = level
}

// logMessageByLevel is the central logging function.
func logMessageByLevel(level LogLevel, message string) {
	mu.Lock()
	defer mu.Unlock()

	// Always write to the file if the file logger is active.
	if fileLoggerActive {
		fileLogger.Printf("[%s] %s", levelToString(level), message)
	}

	// Only write to the console if the message's level is high enough.
	if level >= currentLogLevel {
		log.Printf("[%s] %s", levelToString(level), message)
	}
}

// levelToString converts a LogLevel to a string representation.
func levelToString(level LogLevel) string {
	switch level {
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

// LogDebug writes a message with the DEBUG level.
func LogDebug(message string) {
	logMessageByLevel(DEBUG, message)
}

// LogInfo writes a message with the INFO level.
func LogInfo(message string) {
	logMessageByLevel(INFO, message)
}

// LogWarn writes a message with the WARN level.
func LogWarn(message string) {
	logMessageByLevel(WARN, message)
}

// LogError writes a message with the ERROR level.
func LogError(message string) {
	logMessageByLevel(ERROR, message)
}
