// Package date provides strategies and processors for determining the date of an image.
package date

import "time"

// DateStrategy defines the interface for different methods of obtaining a date from a file.
// Implementations can extract the date from EXIF data, file modification times, etc.
type DateStrategy interface {
	GetDate(fields map[string]interface{}, filePath string) (time.Time, error)
}
