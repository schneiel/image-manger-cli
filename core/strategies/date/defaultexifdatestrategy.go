package date

import (
	"fmt"
	"time"
)

// DefaultExifDateStrategy extracts the date from EXIF metadata.
type DefaultExifDateStrategy struct{}

// NewDefaultExifDateStrategy creates a new DefaultExifDateStrategy.
func NewDefaultExifDateStrategy() *DefaultExifDateStrategy {
	return &DefaultExifDateStrategy{}
}

// Extract extracts date from image metadata fields and file path.
func (s *DefaultExifDateStrategy) Extract(fields map[string]interface{}, _ string) (time.Time, error) {
	if fields == nil {
		return time.Time{}, nil
	}

	dateStr, ok := fields["DateTimeOriginal"].(string)
	if !ok {
		return time.Time{}, nil
	}

	parsedTime, err := time.Parse("2006:01:02 15:04:05", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse EXIF date %s: %w", dateStr, err)
	}
	return parsedTime, nil
}
