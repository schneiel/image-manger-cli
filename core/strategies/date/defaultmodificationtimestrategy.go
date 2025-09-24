// Package date provides different strategies for determining the date of an image.
package date

import (
	"errors"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultModificationTimeStrategy provides functionality for extracting modification time from files.
type DefaultModificationTimeStrategy struct {
	fileSystem filesystem.FileSystem
	localizer  i18n.Localizer
}

// NewDefaultModificationTimeStrategy creates a new DefaultModificationTimeStrategy with injected dependencies.
func NewDefaultModificationTimeStrategy(
	fs filesystem.FileSystem,
	localizer i18n.Localizer,
) (*DefaultModificationTimeStrategy, error) {
	if fs == nil {
		return nil, errors.New("fileSystem cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	return &DefaultModificationTimeStrategy{fileSystem: fs, localizer: localizer}, nil
}

// Extract extracts modification time from the file at the given path.
// Returns the modification time and true if successful, or zero time and false if failed.
func (e *DefaultModificationTimeStrategy) Extract(filePath string, _ map[string]interface{}) (time.Time, bool) {
	fileInfo, err := e.fileSystem.Stat(filePath)
	if err != nil {
		return time.Time{}, false
	}
	return fileInfo.ModTime(), true
}
