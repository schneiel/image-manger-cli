package task

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	"github.com/schneiel/ImageManagerGo/core/processing/dedup"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupaction"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupkeep"
)

// DedupTaskBuilder provides a fluent interface for building DedupTask instances
// with complex dependency injection requirements.
type DedupTaskBuilder struct {
	task *DedupTask
	err  error
}

// NewDedupTaskBuilder creates a new builder instance for constructing DedupTask.
func NewDedupTaskBuilder() *DedupTaskBuilder {
	return &DedupTaskBuilder{
		task: &DedupTask{},
	}
}

// WithConfig sets the configuration for the dedup task.
func (b *DedupTaskBuilder) WithConfig(cfg *config.Config) *DedupTaskBuilder {
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

// WithLogger sets the logger for the dedup task.
func (b *DedupTaskBuilder) WithLogger(logger log.Logger) *DedupTaskBuilder {
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

// WithLocalizer sets the localizer for the dedup task.
func (b *DedupTaskBuilder) WithLocalizer(localizer i18n.Localizer) *DedupTaskBuilder {
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

// WithFilesystem sets the filesystem for the dedup task.
func (b *DedupTaskBuilder) WithFilesystem(fs filesystem.FileSystem) *DedupTaskBuilder {
	if b.err != nil {
		return b
	}
	if fs == nil {
		b.err = errors.New("filesystem cannot be nil")
		return b
	}
	b.task.filesystem = fs
	return b
}

// WithScanner sets the scanner for the dedup task.
func (b *DedupTaskBuilder) WithScanner(scanner dedup.Scanner) *DedupTaskBuilder {
	if b.err != nil {
		return b
	}
	if scanner == nil {
		b.err = errors.New("scanner cannot be nil")
		return b
	}
	b.task.scanner = scanner
	return b
}

// WithHasher sets the hasher for the dedup task.
func (b *DedupTaskBuilder) WithHasher(hasher dedup.Hasher) *DedupTaskBuilder {
	if b.err != nil {
		return b
	}
	if hasher == nil {
		b.err = errors.New("hasher cannot be nil")
		return b
	}
	b.task.hasher = hasher
	return b
}

// WithGrouper sets the grouper for the dedup task.
func (b *DedupTaskBuilder) WithGrouper(grouper dedup.Grouper) *DedupTaskBuilder {
	if b.err != nil {
		return b
	}
	if grouper == nil {
		b.err = errors.New("grouper cannot be nil")
		return b
	}
	if grouper == nil {
		b.err = errors.New("grouper cannot be nil")
		return b
	}
	b.task.grouper = grouper
	return b
}

// WithKeepFunc sets the keep function for the dedup task (new preferred approach).
func (b *DedupTaskBuilder) WithKeepFunc(keepFunc dedupkeep.Func) *DedupTaskBuilder {
	if b.err != nil {
		return b
	}
	if keepFunc == nil {
		b.err = errors.New("keepFunc cannot be nil")
		return b
	}
	b.task.keepFunc = keepFunc
	return b
}

// WithStrategyFactory sets the strategy factory for the dedup task.
func (b *DedupTaskBuilder) WithStrategyFactory(strategyFactory func() dedupaction.Strategy) *DedupTaskBuilder {
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

// WithGroupFlattener sets the group flattener for the dedup task.
func (b *DedupTaskBuilder) WithGroupFlattener(groupFlattener *DefaultGroupFlattener) *DedupTaskBuilder {
	if b.err != nil {
		return b
	}
	b.task.groupFlattener = groupFlattener
	return b
}

// Build constructs the final DedupTask with validation.
func (b *DedupTaskBuilder) Build() (*DedupTask, error) {
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
	if b.task.filesystem == nil {
		return nil, errors.New("filesystem is required")
	}
	if b.task.scanner == nil {
		return nil, errors.New("scanner is required")
	}
	if b.task.hasher == nil {
		return nil, errors.New("hasher is required")
	}
	if b.task.grouper == nil {
		return nil, errors.New("grouper is required")
	}
	if b.task.strategyFactory == nil {
		return nil, errors.New("strategyFactory is required")
	}

	// Validate that keepFunc is provided
	if b.task.keepFunc == nil {
		return nil, errors.New("keepFunc is required")
	}

	// Set default trash path if not provided
	if b.task.config.Deduplicator.TrashPath == "" {
		b.task.config.Deduplicator.TrashPath = filepath.Join(b.task.config.Deduplicator.Source, ".trash")
	}

	// Create group flattener if not provided
	if b.task.groupFlattener == nil {
		groupFlattener, err := NewDefaultGroupFlattener(b.task.logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create group flattener: %w", err)
		}
		b.task.groupFlattener = groupFlattener
	}

	return b.task, nil
}
