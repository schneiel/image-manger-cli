package task

import (
	"errors"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	"github.com/schneiel/ImageManagerGo/core/processing/sort"
)

// SortTaskBuilder provides a fluent interface for building SortTask instances
// with dependency injection requirements.
type SortTaskBuilder struct {
	task *SortTask
	err  error
}

// NewSortTaskBuilder creates a new builder instance for constructing SortTask.
func NewSortTaskBuilder() *SortTaskBuilder {
	return &SortTaskBuilder{
		task: &SortTask{},
	}
}

// WithConfig sets the configuration for the sort task.
func (b *SortTaskBuilder) WithConfig(cfg *config.Config) *SortTaskBuilder {
	if b.err != nil {
		return b
	}
	if cfg == nil {
		b.err = errors.New("config cannot be nil")
		return b
	}
	b.task.config = cfg
	return b
}

// WithLogger sets the logger for the sort task.
func (b *SortTaskBuilder) WithLogger(logger log.Logger) *SortTaskBuilder {
	if b.err != nil {
		return b
	}
	if logger == nil {
		b.err = errors.New("logger cannot be nil")
		return b
	}
	b.task.logger = logger
	return b
}

// WithLocalizer sets the localizer for the sort task.
func (b *SortTaskBuilder) WithLocalizer(localizer i18n.Localizer) *SortTaskBuilder {
	if b.err != nil {
		return b
	}
	if localizer == nil {
		b.err = errors.New("localizer cannot be nil")
		return b
	}
	b.task.localizer = localizer
	return b
}

// WithImageProcessor sets the image processor for the sort task.
func (b *SortTaskBuilder) WithImageProcessor(imageProcessor sort.ImageProcessor) *SortTaskBuilder {
	if b.err != nil {
		return b
	}
	if imageProcessor == nil {
		b.err = errors.New("imageProcessor cannot be nil")
		return b
	}
	b.task.imageProcessor = imageProcessor
	return b
}

// WithFileUtils sets the file utils for the sort task.
func (b *SortTaskBuilder) WithFileUtils(fileUtils filesystem.FileUtils) *SortTaskBuilder {
	if b.err != nil {
		return b
	}
	if fileUtils == nil {
		b.err = errors.New("fileUtils cannot be nil")
		return b
	}
	b.task.fileUtils = fileUtils
	return b
}

// WithStrategyFactory sets the strategy factory for the sort task.
func (b *SortTaskBuilder) WithStrategyFactory(strategyFactory func() sort.Strategy) *SortTaskBuilder {
	if b.err != nil {
		return b
	}
	if strategyFactory == nil {
		b.err = errors.New("strategyFactory cannot be nil")
		return b
	}
	b.task.strategyFactory = strategyFactory
	return b
}

// Build constructs the final SortTask with validation.
func (b *SortTaskBuilder) Build() (*SortTask, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Validate required dependencies
	if b.task.config == nil {
		return nil, errors.New("config is required")
	}
	if b.task.logger == nil {
		return nil, errors.New("logger is required")
	}
	if b.task.localizer == nil {
		return nil, errors.New("localizer is required")
	}
	if b.task.imageProcessor == nil {
		return nil, errors.New("imageProcessor is required")
	}
	if b.task.fileUtils == nil {
		return nil, errors.New("fileUtils is required")
	}
	if b.task.strategyFactory == nil {
		return nil, errors.New("strategyFactory is required")
	}

	return b.task, nil
}
