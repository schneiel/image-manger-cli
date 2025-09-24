// Package date provides interfaces and implementations for extracting dates from files.
package date

import (
	"time"
)

// DateProcessor defines the interface for extracting dates from images
//
//nolint:revive // DateProcessor name is clear and intentional despite package stuttering
type DateProcessor interface {
	GetBestAvailableDate(fields map[string]interface{}, filePath string) (time.Time, error)
}

// Strategy defines the interface for different date extraction strategies.
type Strategy interface {
	GetDate(filePath string) (time.Time, error)
}
