//go:build windows

package date

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// getFileSystemCreationTime retrieves the creation time of a file on Windows systems.
// It uses system-specific calls to access the file's creation time attribute.
func getFileSystemCreationTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	winStat, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}, fmt.Errorf("failed to get Windows file attributes")
	}

	nsec := winStat.CreationTime.Nanoseconds()
	if nsec == 0 {
		return time.Time{}, fmt.Errorf("creation time is zero")
	}

	return time.Unix(0, nsec), nil
}
