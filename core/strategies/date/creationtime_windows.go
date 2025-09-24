//go:build windows

package date

import (
	"errors"
	"syscall"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// GetFileSystemCreationTime retrieves the creation time of a file on Windows systems.
// It uses system-specific calls to access the file's creation time attribute.
func GetFileSystemCreationTime(filePath string, localizer i18n.Localizer) (time.Time, error) {
	return GetFileSystemCreationTimeWithFilesystem(filePath, localizer, &filesystem.DefaultFileSystem{})
}

// GetFileSystemCreationTimeWithFilesystem retrieves the creation time with injected filesystem.
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
		return time.Time{}, err
	}

	winStat, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		// Localized error
		return time.Time{}, errors.New(localizer.Translate("WindowsFileAttributesFailed"))
	}

	nsec := winStat.CreationTime.Nanoseconds()
	if nsec == 0 {
		// Localized error
		return time.Time{}, errors.New(localizer.Translate("CreationTimeIsZero"))
	}

	return time.Unix(0, nsec), nil
}
