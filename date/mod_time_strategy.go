package date

import (
	"os"
	"time"
)

// ModTimeStrategy extracts the date based on the file's last modification time.
type ModTimeStrategy struct{}

// GetDate returns the modification time of the specified file.
func (s *ModTimeStrategy) GetDate(_ map[string]interface{}, filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}
