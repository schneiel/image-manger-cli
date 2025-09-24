// Package shared provides common action strategy patterns for image operations.
package shared

import (
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// DryRunCSVLogger provides common CSV logging functionality for dry-run strategies.
// This utility eliminates duplication between deduplication and sorting dry-run strategies.
type DryRunCSVLogger struct {
	resource  ActionResource
	logger    log.Logger
	localizer i18n.Localizer
}

// NewDryRunCSVLogger creates a new CSV logger for dry-run operations.
func NewDryRunCSVLogger(resource ActionResource, logger log.Logger, localizer i18n.Localizer) *DryRunCSVLogger {
	return &DryRunCSVLogger{
		resource:  resource,
		logger:    logger,
		localizer: localizer,
	}
}

// LogOperation logs a dry-run operation to both console and CSV file.
// operation: The type of operation (e.g., "copy", "move", "delete")
// source: The source file path
// destination: The destination file path (can be empty for delete operations)
func (d *DryRunCSVLogger) LogOperation(operation, source, destination string) error {
	// Log to console if logger is available
	if d.logger != nil {
		if destination != "" {
			d.logger.Infof("[DRY RUN] %s: %s -> %s", operation, source, destination)
		} else {
			d.logger.Infof("[DRY RUN] %s: %s", operation, source)
		}
	}

	// Log to CSV if resource is available
	if d.resource != nil {
		if csvResource, ok := d.resource.(*DefaultCSVResource); ok {
			row := []string{source}
			if destination != "" {
				row = append(row, destination)
			}
			return csvResource.WriteRow(row)
		}
	}

	return nil
}

// LogDeduplicationOperation logs a deduplication-specific dry-run operation.
// original: The original file path
// duplicate: The duplicate file path
func (d *DryRunCSVLogger) LogDeduplicationOperation(original, duplicate string) error {
	d.logger.Infof("[DRY RUN] Potential duplicate found: %s (Original: %s)", duplicate, original)

	if d.resource != nil {
		if csvResource, ok := d.resource.(*DefaultCSVResource); ok {
			return csvResource.WriteRow([]string{original, duplicate})
		}
	}

	return nil
}

// GetResource returns the underlying CSV resource for direct access if needed.
func (d *DryRunCSVLogger) GetResource() ActionResource {
	return d.resource
}
