package di

import (
	"errors"
	"fmt"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// LoggingBuilder handles logging infrastructure setup.
type LoggingBuilder struct{}

// NewLoggingBuilder creates a new logging builder.
func NewLoggingBuilder() *LoggingBuilder {
	return &LoggingBuilder{}
}

// BuildLogging sets up the logging infrastructure.
func (lb *LoggingBuilder) BuildLogging(core *CoreDependencies, cfg *config.Config) (*LoggingDependencies, error) {
	sortLogger, err := lb.createLoggerSafely(core.Localizer, cfg.Sorter.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to create sort logger: %w", err)
	}

	dedupLogger, err := lb.createLoggerSafely(core.Localizer, cfg.Deduplicator.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to create dedup logger: %w", err)
	}

	return &LoggingDependencies{
		SortLogger:  sortLogger,
		DedupLogger: dedupLogger,
	}, nil
}

// createLoggerSafely creates a logger, falling back to console-only if log path is empty.
func (lb *LoggingBuilder) createLoggerSafely(
	localizer i18n.Localizer,
	logPath string,
) (log.Logger, error) { //nolint:ireturn // DI builder returns interface by design
	if logPath == "" {
		// Create console-only logger when log path is empty.
		return log.NewDefaultConsoleLoggerWithWriter(log.INFO, log.NewDefaultConsoleWriter())
	}
	logger, err := log.NewMultiLogger(logPath, localizer)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	return logger, nil
}

// Build is an alias for BuildLogging to satisfy the standard builder pattern.
// This method provides a simplified interface but requires parameters.
func (lb *LoggingBuilder) Build() (*LoggingDependencies, error) {
	// For the standard Build() method, we need to handle the lack of parameters.
	return nil, errors.New("Build() requires parameters - use BuildLogging(core, cfg) instead")
}

// LoggingDependencies holds logging-related dependencies.
type LoggingDependencies struct {
	SortLogger  log.Logger
	DedupLogger log.Logger
}
