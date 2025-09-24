package date

import (
	"errors"
	"time"
)

// ExtractorAdapter adapts an Extractor function to the Strategy interface.
type ExtractorAdapter struct {
	extractor Extractor
}

// NewExtractorAdapter creates a new adapter for the given extractor.
func NewExtractorAdapter(extractor Extractor) Strategy {
	return &ExtractorAdapter{extractor: extractor}
}

// Extract implements the Strategy interface.
func (a *ExtractorAdapter) Extract(fields map[string]interface{}, filePath string) (time.Time, error) {
	if t, ok := a.extractor(filePath, fields); ok {
		return t, nil
	}
	return time.Time{}, errors.New("date not found by extractor")
}
