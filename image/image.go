// Package image defines the core data structure for an image file.
package image

import "time"

// Image represents an image file with its most important metadata.
type Image struct {
	FilePath         string
	OriginalFileName string
	Date             time.Time
}
