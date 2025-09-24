// Package log provides SOLID-compliant logging components.
package log

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultFileLogger logs messages to a file.
type DefaultFileLogger struct {
	logger     *log.Logger
	logFile    *os.File
	minLevel   Level
	localizer  i18n.Localizer
	fileSystem filesystem.FileSystem
}

// NewDefaultFileLogger creates a new logger that writes to the specified file.
// It logs all messages regardless of the level set elsewhere, as in the original code.
func NewDefaultFileLogger(filePath string, minLevel Level) (*DefaultFileLogger, error) {
	// #nosec G304 -- filePath is controlled by application configuration
	// #nosec G302 -- 0o600 is appropriate for log files to restrict access
	return NewDefaultFileLoggerWithDependencies(filePath, minLevel, nil, &filesystem.DefaultFileSystem{})
}

// NewDefaultFileLoggerWithLocalizer creates a new logger with a provided localizer.
func NewDefaultFileLoggerWithLocalizer(
	filePath string,
	minLevel Level,
	localizer i18n.Localizer,
) (*DefaultFileLogger, error) {
	return NewDefaultFileLoggerWithDependencies(filePath, minLevel, localizer, &filesystem.DefaultFileSystem{})
}

// NewDefaultFileLoggerWithDependencies creates a new logger with injected dependencies.
func NewDefaultFileLoggerWithDependencies(
	filePath string,
	minLevel Level,
	localizer i18n.Localizer,
	fileSystem filesystem.FileSystem,
) (*DefaultFileLogger, error) {
	if fileSystem == nil {
		return nil, errors.New("fileSystem cannot be nil")
	}

	file, err := fileSystem.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		if localizer != nil {
			return nil, fmt.Errorf(
				localizer.Translate("ErrorOpeningLogFile", map[string]interface{}{"Path": filePath, "Error": err}),
				err,
			)
		}
		return nil, fmt.Errorf("ErrorOpeningLogFile: %w", err)
	}

	return &DefaultFileLogger{
		logger:     log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile),
		logFile:    file.(*os.File), // Type assertion for backward compatibility
		minLevel:   minLevel,
		localizer:  localizer,
		fileSystem: fileSystem,
	}, nil
}

// Close closes the log file.
func (l *DefaultFileLogger) Close() error {
	if l.logFile != nil {
		err := l.logFile.Close()
		l.logFile = nil // Set to nil to prevent double close
		if err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}
	return nil
}

// SetLevel configures the minimum log level for file output. Messages below this level are filtered out.
// File logs typically use DEBUG or INFO to capture detailed operation history for troubleshooting.
func (l *DefaultFileLogger) SetLevel(level Level) {
	l.minLevel = level
}

func (l *DefaultFileLogger) log(level Level, format string, v ...interface{}) {
	if level < l.minLevel {
		return
	}
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[%s] %s\n", level.String(), message)
}

// Public log methods

// Debug logs a debug level message.
func (l *DefaultFileLogger) Debug(message string) { l.log(DEBUG, "%s", message) }

// Info logs an info level message.
func (l *DefaultFileLogger) Info(message string) { l.log(INFO, "%s", message) }

// Warn logs a warning level message.
func (l *DefaultFileLogger) Warn(message string) { l.log(WARN, "%s", message) }

func (l *DefaultFileLogger) Error(message string) { l.log(ERROR, "%s", message) }

// Public formatted log methods

// Debugf logs a formatted debug level message.
func (l *DefaultFileLogger) Debugf(format string, v ...interface{}) { l.log(DEBUG, format, v...) }

// Infof logs a formatted info level message.
func (l *DefaultFileLogger) Infof(format string, v ...interface{}) { l.log(INFO, format, v...) }

// Warnf logs a formatted warning level message.
func (l *DefaultFileLogger) Warnf(format string, v ...interface{}) { l.log(WARN, format, v...) }

// Errorf logs a formatted error level message.
func (l *DefaultFileLogger) Errorf(format string, v ...interface{}) { l.log(ERROR, format, v...) }
