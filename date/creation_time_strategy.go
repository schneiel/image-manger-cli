package date

import "time"

// CreationTimeStrategy extracts the date based on the file's creation time from the file system.
type CreationTimeStrategy struct{}

// GetDate returns the creation time of the specified file.
func (s *CreationTimeStrategy) GetDate(_ map[string]interface{}, filePath string) (time.Time, error) {
	return getFileSystemCreationTime(filePath)
}
