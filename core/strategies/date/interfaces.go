// Package date provides different strategies for determining the date of an image.
package date

import "time"

// Strategy defines the interface for extracting dates from image metadata and files.
// Implementing classes should provide specific logic for different date extraction methods.
type Strategy interface {
	// Extract extracts date from image metadata fields and file path
	Extract(fields map[string]interface{}, filePath string) (time.Time, error)
}

// Extractor is a function type that defines a strategy for extracting a date.
// It returns the extracted time and a boolean indicating success.
type Extractor func(_ string, _ map[string]interface{}) (time.Time, bool)
