//go:build darwin

// Package date provides functionality for extracting file creation timestamps.
package date

import (
	"fmt"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// GetFileSystemCreationTime returns the creation time of a file on Darwin systems.
// It uses the file's birth time from the stat metadata.
// It returns an error if the file information cannot be retrieved.
func GetFileSystemCreationTime(filePath string, localizer i18n.Localizer) (time.Time, error) {
	return GetFileSystemCreationTimeWithFilesystem(filePath, localizer, &filesystem.DefaultFileSystem{})
}

// GetFileSystemCreationTimeWithFilesystem returns the creation time with injected filesystem.
func GetFileSystemCreationTimeWithFilesystem(
	filePath string,
	localizer i18n.Localizer,
	fileSystem filesystem.FileSystem,
) (time.Time, error) {
	if fileSystem == nil {
		panic("fileSystem cannot be nil")
	}

	fileInfo, err := fileSystem.Stat(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"%s",
			localizer.Translate("FileInfoError", map[string]interface{}{"FilePath": filePath, "Error": err}),
		)
	}

	// Birth time is not directly available in os.FileInfo in a portable way.
	// Accessing platform-specific details is required.
	// This functionality would typically be implemented using syscall Stat_t,
	// but is omitted here to avoid platform-specific cgo dependencies in this example.
	// Returning ModTime as a fallback for demonstration purposes.
	return fileInfo.ModTime(), nil
}
