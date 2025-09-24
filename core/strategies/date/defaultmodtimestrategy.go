package date

import (
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
)

// DefaultModTimeStrategy extracts the modification time from file info.
type DefaultModTimeStrategy struct {
	fileSystem filesystem.FileSystem
}

// NewDefaultModTimeStrategy creates a new DefaultModTimeStrategy.
func NewDefaultModTimeStrategy() *DefaultModTimeStrategy {
	return NewDefaultModTimeStrategyWithFilesystem(&filesystem.DefaultFileSystem{})
}

// NewDefaultModTimeStrategyWithFilesystem creates a new DefaultModTimeStrategy with injected filesystem.
func NewDefaultModTimeStrategyWithFilesystem(fileSystem filesystem.FileSystem) *DefaultModTimeStrategy {
	if fileSystem == nil {
		panic("fileSystem cannot be nil")
	}
	return &DefaultModTimeStrategy{fileSystem: fileSystem}
}

// Extract extracts date from image metadata fields and file path.
func (s *DefaultModTimeStrategy) Extract(_ map[string]interface{}, filePath string) (time.Time, error) {
	fileInfo, err := s.fileSystem.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}
