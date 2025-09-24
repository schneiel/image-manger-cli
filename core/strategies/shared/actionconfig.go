// Package shared provides common action strategy patterns for image operations.
package shared

import (
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
)

// ActionConfig holds the dependencies for action strategies across all domains.
// This unified configuration supports both image deduplication and sorting strategies.
type ActionConfig struct {
	Logger    log.Logger
	Localizer i18n.Localizer
	FileUtils filesystem.FileUtils
	CsvPath   string
	TrashPath string // Optional - used only by deduplication strategies that move files to trash
}
