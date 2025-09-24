// Package date provides functionality for extracting file creation times.
package date

import (
	"errors"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// DefaultCreationTimeStrategy provides functionality for extracting file creation times.
type DefaultCreationTimeStrategy struct {
	fileSystem filesystem.FileSystem
	localizer  i18n.Localizer
}

// NewDefaultCreationTimeStrategy creates a new DefaultCreationTimeStrategy with injected dependencies.
func NewDefaultCreationTimeStrategy(fs filesystem.FileSystem, localizer i18n.Localizer) (*DefaultCreationTimeStrategy, error) {
	if fs == nil {
		return nil, errors.New("fileSystem cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	return &DefaultCreationTimeStrategy{fileSystem: fs, localizer: localizer}, nil
}

// Extract extracts the creation time from the file at the given path.
// Returns the creation time or an error if extraction fails.
func (e *DefaultCreationTimeStrategy) Extract(filePath string) (time.Time, error) {
	return e.extractPlatformSpecific(filePath)
}

// extractPlatformSpecific contains the platform-specific implementation.
// This method will be implemented differently for each platform.
func (e *DefaultCreationTimeStrategy) extractPlatformSpecific(filePath string) (time.Time, error) {
	// This will be implemented by platform-specific files
	// For now, return modification time as fallback
	fileInfo, err := e.fileSystem.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}
