//go:build darwin

package date

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// getFileSystemCreationTime retrieves the creation time of a file on Darwin-based systems.
// It uses system-specific calls to access the birth time of the file.
func getFileSystemCreationTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, fmt.Errorf("failed to get Unix file stats")
	}
	return time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec), nil
}
