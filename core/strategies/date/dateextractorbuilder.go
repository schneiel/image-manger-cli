package date

import (
	"errors"
	"fmt"
	"time"

	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
)

// ExtractorBuilder provides a fluent interface for building date extractor chains.
type ExtractorBuilder struct {
	extractors []Extractor
	localizer  i18n.Localizer
	filesystem filesystem.FileSystem
	err        error
}

// NewExtractorBuilder creates a new builder instance for constructing date extractors.
func NewExtractorBuilder(localizer i18n.Localizer) *ExtractorBuilder {
	if localizer == nil {
		return &ExtractorBuilder{err: errors.New("localizer cannot be nil")}
	}
	return &ExtractorBuilder{
		extractors: make([]Extractor, 0),
		localizer:  localizer,
	}
}

// WithFilesystem sets the filesystem for operations that require file access.
func (b *ExtractorBuilder) WithFilesystem(fs filesystem.FileSystem) *ExtractorBuilder {
	if b.err != nil {
		return b
	}
	if fs == nil {
		b.err = errors.New("filesystem cannot be nil")
		return b
	}
	b.filesystem = fs
	return b
}

// AddExifExtractor adds an EXIF-based date extractor with specified field and layout.
func (b *ExtractorBuilder) AddExifExtractor(fieldName, layout string) *ExtractorBuilder {
	if b.err != nil {
		return b
	}
	if fieldName == "" {
		b.err = errors.New("fieldName cannot be empty")
		return b
	}
	if layout == "" {
		b.err = errors.New("layout cannot be empty")
		return b
	}

	extractor := func(_ string, exifData map[string]interface{}) (time.Time, bool) {
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

	b.extractors = append(b.extractors, extractor)
	return b
}

// AddModTimeExtractor adds a modification time extractor.
func (b *ExtractorBuilder) AddModTimeExtractor() *ExtractorBuilder {
	if b.err != nil {
		return b
	}
	if b.filesystem == nil {
		b.err = errors.New("filesystem is required for modTime extractor")
		return b
	}

	modTimeExtractor, err := NewDefaultModificationTimeStrategy(b.filesystem, b.localizer)
	if err != nil {
		b.err = err
		return b
	}
	b.extractors = append(b.extractors, modTimeExtractor.Extract)
	return b
}

// AddCreationTimeExtractor adds a creation time extractor.
func (b *ExtractorBuilder) AddCreationTimeExtractor() *ExtractorBuilder {
	if b.err != nil {
		return b
	}
	if b.filesystem == nil {
		b.err = errors.New("filesystem is required for creationTime extractor")
		return b
	}

	extractor, err := NewDefaultCreationTimeStrategy(b.filesystem, b.localizer)
	if err != nil {
		b.err = err
		return b
	}
	creationTimeExtractor := func(filePath string, _ map[string]interface{}) (time.Time, bool) {
		t, err := extractor.Extract(filePath)
		if err != nil {
			return time.Time{}, false
		}
		return t, true
	}

	b.extractors = append(b.extractors, creationTimeExtractor)
	return b
}

// AddCustomExtractor adds a custom date extractor.
func (b *ExtractorBuilder) AddCustomExtractor(extractor Extractor) *ExtractorBuilder {
	if b.err != nil {
		return b
	}
	if extractor == nil {
		b.err = errors.New("extractor cannot be nil")
		return b
	}

	b.extractors = append(b.extractors, extractor)
	return b
}

// AddExifExtractors adds multiple EXIF extractors from configuration.
func (b *ExtractorBuilder) AddExifExtractors(exifConfigs []ExifStrategyConfig) *ExtractorBuilder {
	if b.err != nil {
		return b
	}

	for _, config := range exifConfigs {
		b.AddExifExtractor(config.FieldName, config.Layout)
		if b.err != nil {
			return b
		}
	}

	return b
}

// BuildChain creates a single Extractor that tries all added extractors in order.
func (b *ExtractorBuilder) BuildChain() (Extractor, error) {
	if b.err != nil {
		return nil, b.err
	}

	if len(b.extractors) == 0 {
		return nil, errors.New("at least one extractor must be added")
	}

	extractors := make([]Extractor, len(b.extractors))
	copy(extractors, b.extractors)

	return func(filePath string, exifData map[string]interface{}) (time.Time, bool) {
		for _, extractor := range extractors {
			if t, ok := extractor(filePath, exifData); ok {
				return t, true
			}
		}
		return time.Time{}, false
	}, nil
}

// BuildList returns a slice of all added extractors.
func (b *ExtractorBuilder) BuildList() ([]Extractor, error) {
	if b.err != nil {
		return nil, b.err
	}

	if len(b.extractors) == 0 {
		return nil, errors.New("at least one extractor must be added")
	}

	extractors := make([]Extractor, len(b.extractors))
	copy(extractors, b.extractors)
	return extractors, nil
}

// Build creates a single Extractor that tries all added extractors in order
// This is an alias for BuildChain() to satisfy the standard builder pattern.
func (b *ExtractorBuilder) Build() (Extractor, error) {
	return b.BuildChain()
}

// BuildFromConfig creates extractors based on configuration (compatible with existing factory).
func (b *ExtractorBuilder) BuildFromConfig(config ExtractorConfig) ([]Extractor, error) {
	if b.err != nil {
		return nil, b.err
	}
	if config.Localizer == nil {
		return nil, errors.New("localizer is required")
	}
	if config.FileSystem == nil {
		return nil, errors.New("fileSystem is required")
	}

	// Update builder with config dependencies
	b.localizer = config.Localizer
	b.filesystem = config.FileSystem

	for _, strategyName := range config.StrategyOrder {
		switch strategyName {
		case "exif":
			b.AddExifExtractors(config.ExifStrategies)
		case "modTime":
			b.AddModTimeExtractor()
		case "creationTime":
			b.AddCreationTimeExtractor()
		default:
			return nil, fmt.Errorf("unknown date strategy: %s", strategyName)
		}

		if b.err != nil {
			return nil, b.err
		}
	}

	return b.BuildList()
}
