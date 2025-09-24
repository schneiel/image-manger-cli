package log

import (
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// NewMultiLogger creates a multi-logger with file and console output.
// This replaces the over-engineered LoggerFactory pattern with a simple constructor.
func NewMultiLogger(logFile string, localizer i18n.Localizer) (Logger, error) {
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}

	fileLogger, err := NewDefaultFileLoggerWithDependencies(logFile, WARN, localizer, &filesystem.DefaultFileSystem{})
	if err != nil {
		return nil, fmt.Errorf(
			"%s",
			localizer.Translate("FileLoggerCreationError", map[string]interface{}{"Error": err}),
		)
	}

	consoleLogger, err := NewDefaultConsoleLoggerWithWriter(INFO, NewDefaultConsoleWriter())
	if err != nil {
		return nil, fmt.Errorf("failed to create console logger: %w", err)
	}
	return NewDefaultMultiLogger(consoleLogger, fileLogger), nil
}
