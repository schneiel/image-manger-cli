package task

import (
	"errors"
	"path/filepath"

	"github.com/schneiel/ImageManagerGo/core/config"
	"github.com/schneiel/ImageManagerGo/core/image"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/filesystem"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/i18n"
	"github.com/schneiel/ImageManagerGo/core/infrastructure/log"
	"github.com/schneiel/ImageManagerGo/core/processing/dedup"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupaction"
	"github.com/schneiel/ImageManagerGo/core/strategies/dedupkeep"
)

// DedupTask implements the Task interface for image deduplication.
type DedupTask struct {
	config     *config.Config
	logger     log.Logger
	localizer  i18n.Localizer
	filesystem filesystem.FileSystem
	scanner    dedup.Scanner
	hasher     dedup.Hasher
	grouper    dedup.Grouper

	keepFunc       dedupkeep.Func         // Function-based keep logic
	groupFlattener *DefaultGroupFlattener // Group flattening utility
	
	// Lazy strategy creation
	strategyFactory func() dedupaction.Strategy
	cachedStrategy  dedupaction.Strategy
	lastConfigValue string // Track config changes
}

// NewDedupTask creates a new dedup task using function-based keep logic.
func NewDedupTask(
	cfg *config.Config,
	logger log.Logger,
	localizer i18n.Localizer,
	fsys filesystem.FileSystem,
	scanner dedup.Scanner,
	hasher dedup.Hasher,
	grouper dedup.Grouper,
	keepFunc dedupkeep.Func, // Function-based keep logic
	strategyFactory func() dedupaction.Strategy, // Strategy factory for lazy creation
) (*DedupTask, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if localizer == nil {
		return nil, errors.New("localizer cannot be nil")
	}
	if fsys == nil {
		return nil, errors.New("filesystem cannot be nil")
	}
	if scanner == nil {
		return nil, errors.New("scanner cannot be nil")
	}
	if hasher == nil {
		return nil, errors.New("hasher cannot be nil")
	}
	if grouper == nil {
		return nil, errors.New("grouper cannot be nil")
	}
	if keepFunc == nil {
		return nil, errors.New("keepFunc cannot be nil")
	}
	if strategyFactory == nil {
		return nil, errors.New("strategyFactory cannot be nil")
	}

	// Set default trash path if not provided
	if cfg.Deduplicator.TrashPath == "" {
		cfg.Deduplicator.TrashPath = filepath.Join(cfg.Deduplicator.Source, ".trash")
	}

	// Create group flattener for utility operations
	groupFlattener, err := NewDefaultGroupFlattener(logger)
	if err != nil {
		return nil, err
	}

	return &DedupTask{
		config:          cfg,
		logger:          logger,
		localizer:       localizer,
		filesystem:      fsys,
		scanner:         scanner,
		hasher:          hasher,
		grouper:         grouper,
		keepFunc:        keepFunc,
		groupFlattener:  groupFlattener,
		strategyFactory: strategyFactory,
	}, nil
}

// getStrategy returns the action strategy, creating it lazily or rebuilding if config changed
func (t *DedupTask) getStrategy() dedupaction.Strategy {
	currentConfigValue := t.config.Deduplicator.ActionStrategy
	
	// Rebuild strategy if config changed or never built
	if t.cachedStrategy == nil || t.lastConfigValue != currentConfigValue {
		t.cachedStrategy = t.strategyFactory()
		t.lastConfigValue = currentConfigValue
	}
	
	return t.cachedStrategy
}

// Run executes the dedup task using injected dependencies.
func (t *DedupTask) Run() error {
	// Get strategy (lazy creation with config change detection)
	actionStrategy := t.getStrategy()
	
	if actionStrategy.GetResources() != nil {
		err := actionStrategy.GetResources().Setup()
		if err != nil {
			// Handle nil localizer gracefully (e.g., in tests)
			errorMsg := "ActionStrategyError"
			if t.localizer != nil {
				errorMsg = t.localizer.Translate("ActionStrategyError", map[string]interface{}{"Error": err})
			}
			return errors.New(errorMsg + ": " + err.Error())
		}
		defer func() { _ = actionStrategy.GetResources().Teardown() }()
	}

	potentialGroups, err := t.scanner.Scan(t.config.Deduplicator.Source)
	if err != nil {
		return err
	}

	filesToHash := t.groupFlattener.Flatten(potentialGroups)
	hashes, err := t.hasher.HashFiles(filesToHash)
	if err != nil {
		return err
	}

	duplicateGroups, err := t.grouper.Group(hashes)
	if err != nil {
		return err
	}

	t.processGroups(duplicateGroups, t.keepFunc, actionStrategy)
	return nil
}

func (t *DedupTask) processGroups(
	groups []dedup.DuplicateGroup,
	keepFunc dedupkeep.Func,
	actionStrategy dedupaction.Strategy,
) {
	if len(groups) == 0 {
		msg := "SummaryNoDuplicates"
		if t.localizer != nil {
			msg = t.localizer.Translate("SummaryNoDuplicates")
		}
		t.logger.Info(msg)
		return
	}

	filesToRemoveCount := 0
	for _, group := range groups {
		// Convert group to string slice for compatibility
		groupFiles := []string(group)

		// Use the keep function to determine which file to keep
		toKeep, toRemove := keepFunc(groupFiles)

		filesToRemoveCount += len(toRemove)

		msg := "DuplicateGroupFound"
		if t.localizer != nil {
			msg = t.localizer.Translate(
				"DuplicateGroupFound",
				map[string]interface{}{"ToKeep": toKeep, "ToRemoveCount": len(toRemove)},
			)
		}
		t.logger.Info(msg)

		// Execute action using the strategy pattern
		toKeepImage := &image.Image{FilePath: toKeep}
		for _, toRemovePath := range toRemove {
			toRemoveImage := &image.Image{FilePath: toRemovePath}
			err := actionStrategy.Execute(toKeepImage, toRemoveImage)
			if err != nil {
				t.logger.Errorf(
					"Failed to execute action on group (keep: %s, remove: %s): %v",
					toKeep,
					toRemovePath,
					err,
				)
			}
		}
	}

	summaryMsg := "SummaryDuplicatesFound"
	if t.localizer != nil {
		summaryMsg = t.localizer.Translate(
			"SummaryDuplicatesFound",
			map[string]interface{}{"Groups": len(groups), "Files": filesToRemoveCount},
		)
	}
	t.logger.Info(summaryMsg)
}

// flatten function has been replaced by DefaultGroupFlattener.Flatten()
