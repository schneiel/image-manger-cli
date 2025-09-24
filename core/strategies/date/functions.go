// Package date provides standalone functions for creating date extractors.
package date

import (
	"errors"
	"fmt"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// ExtractorConfig holds configuration for date extractors.
type ExtractorConfig struct {
	StrategyOrder  []string
	ExifStrategies []ExifStrategyConfig
	Localizer      i18n.Localizer
	FileSystem     filesystem.FileSystem
}

// ExifStrategyConfig holds configuration for EXIF-based extractors.
type ExifStrategyConfig struct {
	FieldName string
	Layout    string
}

// NewExifExtractor creates an EXIF-based date extractor.
// This replaces the over-engineered factory pattern with a simple constructor.
func NewExifExtractor(fieldName, layout string) Extractor {
	return func(_ string, exifData map[string]interface{}) (time.Time, bool) {
		dateValue, ok := exifData[fieldName]
		if !ok {
			return time.Time{}, false
		}

		dateString, ok := dateValue.(string)
		if !ok {
			return time.Time{}, false
		}

		t, err := time.Parse(layout, dateString)
		if err != nil {
			return time.Time{}, false
		}

		return t, true
	}
}

// NewModTimeExtractor creates a modification time extractor.
// This replaces the over-engineered factory pattern with a simple constructor.
func NewModTimeExtractor(fileSystem filesystem.FileSystem, localizer i18n.Localizer) (Extractor, error) {
	modTimeExtractor, err := NewDefaultModificationTimeStrategy(fileSystem, localizer)
	if err != nil {
		return nil, err
	}
	return modTimeExtractor.Extract, nil
}

// NewCreationTimeExtractor creates a creation time extractor.
// This replaces the over-engineered factory pattern with a simple constructor.
func NewCreationTimeExtractor(fileSystem filesystem.FileSystem, localizer i18n.Localizer) (Extractor, error) {
	extractor, err := NewDefaultCreationTimeStrategy(fileSystem, localizer)
	if err != nil {
		return nil, err
	}
	return func(filePath string, _ map[string]interface{}) (time.Time, bool) {
		t, err := extractor.Extract(filePath)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	}, nil
}

// CreateExtractorsFromConfig creates extractors based on configuration.
func CreateExtractorsFromConfig(config ExtractorConfig) ([]Extractor, error) {
	if config.Localizer == nil {
		return nil, errors.New("localizer is required")
	}
	if config.FileSystem == nil {
		return nil, errors.New("fileSystem is required")
	}

	var extractors []Extractor

	for _, strategyName := range config.StrategyOrder {
		switch strategyName {
		case "exif":
			extractors = append(extractors, createExifExtractors(config.ExifStrategies, config.Localizer)...)
		case "modTime":
			extractor, err := NewModTimeExtractor(config.FileSystem, config.Localizer)
			if err != nil {
				return nil, err
			}
			extractors = append(extractors, extractor)
		case "creationTime":
			extractor, err := NewCreationTimeExtractor(config.FileSystem, config.Localizer)
			if err != nil {
				return nil, err
			}
			extractors = append(extractors, extractor)
		default:
			return nil, fmt.Errorf("unknown date strategy: %s", strategyName)
		}
	}

	return extractors, nil
}

// NewExtractorChain creates a single Extractor that tries multiple extractors in order.
// It returns the result of the first successful extractor.
// This replaces the over-engineered factory pattern with a simple constructor.
func NewExtractorChain(extractors ...Extractor) Extractor {
	return func(filePath string, exifData map[string]interface{}) (time.Time, bool) {
		for _, extractor := range extractors {
			if t, ok := extractor(filePath, exifData); ok {
				return t, true
			}
		}
		return time.Time{}, false
	}
}

// createExifExtractors creates EXIF-based date extractors (helper function).
func createExifExtractors(exifConfigs []ExifStrategyConfig, _ i18n.Localizer) []Extractor {
	extractors := make([]Extractor, 0, len(exifConfigs))
	for _, exifCfg := range exifConfigs {
		extractors = append(extractors, NewExifExtractor(exifCfg.FieldName, exifCfg.Layout))
	}
	return extractors
}
