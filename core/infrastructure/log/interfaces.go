package log

// Logger defines the interface for all loggers.
type Logger interface {
	SetLevel(level Level)
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}
