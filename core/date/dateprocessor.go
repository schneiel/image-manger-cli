// Package date provides functionality for extracting dates from images.
package date

import (
	"fmt"
	"time"

	datestrategies "github.com/schneiel/ImageManagerGo/core/strategies/date"
)

// FieldName defines a type for EXIF field names for type safety.
type FieldName string

// StrategyProcessor uses the OOP strategy pattern for date extraction.
type StrategyProcessor struct {
	strategy datestrategies.Strategy
}

// NewStrategyProcessor creates a new processor using the strategy pattern.
func NewStrategyProcessor(dateStrategy datestrategies.Strategy) *StrategyProcessor {
	return &StrategyProcessor{strategy: dateStrategy}
}

// GetBestAvailableDate uses the strategy to find the best available date.
func (sp *StrategyProcessor) GetBestAvailableDate(fields map[string]interface{}, filePath string) (time.Time, error) {
	date, err := sp.strategy.Extract(fields, filePath)
	if err == nil && !date.IsZero() {
		return date.Local(), nil
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to extract date from %s: %w", filePath, err)
	}
	return time.Time{}, nil
}
