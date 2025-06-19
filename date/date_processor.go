package date

import (
	"ImageManager/i18n"
	"fmt"
	"time"
)

// FieldName defines a type for EXIF field names for type safety.
type FieldName string

// Constants for EXIF date-related field names.
const (
	DateTimeOriginal       FieldName = "DateTimeOriginal"
	DateTimeDigitized      FieldName = "DateTimeDigitized"
	SubSecDateTimeOriginal FieldName = "SubSecDateTimeOriginal"
	DateTime               FieldName = "DateTime"
)

// DatePriorityProcessor determines the best available date for a file by trying a series of strategies in order.
type DatePriorityProcessor struct {
	strategies []DateStrategy
}

// NewDatePriorityProcessor creates a new processor with a predefined list of date-finding strategies.
// The strategies are ordered from most to least preferred.
func NewDatePriorityProcessor() *DatePriorityProcessor {
	// This now represents the default order
	processor, _ := NewDatePriorityProcessorFromStrategies([]string{"exif", "modTime", "creationTime"})
	return processor
}

// NewDatePriorityProcessorFromStrategies creates a new processor from a list of strategy names.
func NewDatePriorityProcessorFromStrategies(strategyNames []string) (*DatePriorityProcessor, error) {
	var strategies []DateStrategy
	for _, name := range strategyNames {
		switch name {
		case "exif":
			strategies = append(strategies,
				&ExifStrategy{fieldName: DateTimeOriginal, layout: "2006:01:02 15:04:05"},
				&ExifStrategy{fieldName: DateTimeDigitized, layout: "2006:01:02 15:04:05"},
				&ExifStrategy{fieldName: SubSecDateTimeOriginal, layout: "2006:01:02 15:04:05.00"},
				&ExifStrategy{fieldName: DateTime, layout: "2006:01:02 15:04:05-07:00"},
			)
		case "modTime":
			strategies = append(strategies, &ModTimeStrategy{})
		case "creationTime":
			strategies = append(strategies, &CreationTimeStrategy{})
		default:
			return nil, fmt.Errorf("unknown date strategy: %s", name)
		}
	}
	return &DatePriorityProcessor{strategies: strategies}, nil
}

// GetBestAvailableDate iterates through its strategies and returns the first valid date found.
func (dpp *DatePriorityProcessor) GetBestAvailableDate(fields map[string]interface{}, filePath string) (time.Time, error) {
	for _, strategy := range dpp.strategies {
		date, err := strategy.GetDate(fields, filePath)
		if err == nil && !date.IsZero() {
			return date.Local(), nil
		}
	}

	return time.Time{}, fmt.Errorf(i18n.T("ErrorNoValidDate", map[string]interface{}{"Path": filePath}))
}
