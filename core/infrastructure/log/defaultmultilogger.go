// Package log provides SOLID-compliant logging components.
package log

// DefaultMultiLogger distributes log messages to multiple loggers.
type DefaultMultiLogger struct {
	loggers []Logger
}

// NewDefaultMultiLogger creates a new multi-logger with the given loggers.
func NewDefaultMultiLogger(loggers ...Logger) *DefaultMultiLogger {
	return &DefaultMultiLogger{
		loggers: loggers,
	}
}

// SetLevel sets the level for all contained loggers.
func (ml *DefaultMultiLogger) SetLevel(level Level) {
	for _, logger := range ml.loggers {
		logger.SetLevel(level)
	}
}

// Debug sends a debug message to all loggers.
func (ml *DefaultMultiLogger) Debug(message string) {
	for _, logger := range ml.loggers {
		logger.Debug(message)
	}
}

// Info sends an info message to all loggers.
func (ml *DefaultMultiLogger) Info(message string) {
	for _, logger := range ml.loggers {
		logger.Info(message)
	}
}

// Warn sends a warning message to all loggers.
func (ml *DefaultMultiLogger) Warn(message string) {
	for _, logger := range ml.loggers {
		logger.Warn(message)
	}
}

// Error sends an error message to all loggers.
func (ml *DefaultMultiLogger) Error(message string) {
	for _, logger := range ml.loggers {
		logger.Error(message)
	}
}

// Debugf sends a formatted debug message to all loggers.
func (ml *DefaultMultiLogger) Debugf(format string, v ...interface{}) {
	for _, logger := range ml.loggers {
		logger.Debugf(format, v...)
	}
}

// Infof sends a formatted info message to all loggers.
func (ml *DefaultMultiLogger) Infof(format string, v ...interface{}) {
	for _, logger := range ml.loggers {
		logger.Infof(format, v...)
	}
}

// Warnf sends a formatted warning message to all loggers.
func (ml *DefaultMultiLogger) Warnf(format string, v ...interface{}) {
	for _, logger := range ml.loggers {
		logger.Warnf(format, v...)
	}
}

// Errorf sends a formatted error message to all loggers.
func (ml *DefaultMultiLogger) Errorf(format string, v ...interface{}) {
	for _, logger := range ml.loggers {
		logger.Errorf(format, v...)
	}
}
